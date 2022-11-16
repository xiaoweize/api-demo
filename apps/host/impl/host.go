package impl

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/xiaoweize/api-demo/apps/host"
)

//业务处理层(Controller层 )
//结构体HostServiceImpl方法 实现了host.Service接口
func (i *HostServiceImpl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	i.l.Info("create host")
	//Named创建子logger 分别创建create/test子logger对应的debug/info级别日志
	i.l.Named("Create").Debug("create host")
	i.l.Named("test").Info("create host")
	//格式化打印debug日志
	i.l.Debugf("create host %s", ins.Name)
	//携带额外的meta数据，常用于Trace系统
	i.l.With(logger.NewAny("request-id", "req01")).Debug("create host with meta kv")

	// 校验数据合法性
	if err := ins.Validate(); err != nil {
		return nil, err
	}

	//默认值字段填充
	ins.InjectDefault()
	//后面是与mysql服务器的交互 放在dao层
	err := i.save(ctx, ins)
	if err != nil {
		return nil, err
	}
	return ins, nil
}

//查询主机列表
func (i *HostServiceImpl) QueryHost(ctx context.Context, req *host.QueryHostRequest) (*host.HostSet, error) {
	b := sqlbuilder.NewBuilder(QueryHostSql) //基础语句
	//如果请求中带了查询关键字就带上就添加where语句
	if req.Keywords != "" {
		//在基础语句上添加where条件
		b.Where("r.`name` like ? or r.description like ? or r.private_ip like ? or r.public_ip like ?",
			"%"+req.Keywords+"%",
			"%"+req.Keywords+"%",
			//IP地址走后缀匹配  索引效率更高
			req.Keywords+"%",
			req.Keywords+"%",
		)
	}
	//添加limit这是客户端requst请求传进来的值
	b.Limit(req.OffSet(), uint(req.PageSize))
	//生成sql语句和参数
	querySql, args := b.Build()
	fmt.Println(args...)
	i.l.Debugf("query sql:%v,args:%v", querySql, args)

	//后面是与mysql服务器的交互
	//构建一个prepare statment prepare会发送给数据库服务器端
	stmt, err := i.db.PrepareContext(ctx, querySql)
	if err != nil {
		return nil, err
	}
	//注意一定要关闭statment
	defer stmt.Close()
	//将args参数发送给mysql服务器端执行,返回结果rows
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	//rows资源也要延迟关闭
	defer rows.Close()

	set := host.NewHostSet()
	for rows.Next() {
		//将扫描的数据读取出来转换到对象中,要获取如下字段
		//h.cpu,h.memory,h.gpu_spec,h.gpu_amount,h.os_type,h.os_name,h.serial_number
		ins := host.NewHost()
		if err := rows.Scan(&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt, &ins.ExpireAt,
			&ins.Type, &ins.Name, &ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt,
			&ins.Account, &ins.PublicIP, &ins.PrivateIP,
			&ins.CPU, &ins.Memory, &ins.GPUSpec, &ins.GPUAmount, &ins.OSType, &ins.OSName, &ins.SerialNumber); err != nil {
			return nil, err
		}
		set.Add(ins)
	}

	// total统计关键字查询结果数量,注意BuildCount代码里面去掉了上面的limit语句，这样客户端发送请求时如果不带关键字查询时total显示的是所有主机数
	//这也符合用户逻辑：如阿里云主机页面，不输入关键字查询会显示所有主机数，但分页功能正常
	countSQL, args := b.BuildCount()
	i.l.Debugf("count sql: %s, args: %v", countSQL, args)
	countStmt, err := i.db.PrepareContext(ctx, countSQL)
	if err != nil {
		return nil, err
	}
	defer countStmt.Close()
	if err := countStmt.QueryRowContext(ctx, args...).Scan(&set.Total); err != nil {
		return nil, err
	}

	return set, nil

}

//查询主机详情
func (i *HostServiceImpl) DescribeHost(ctx context.Context, req *host.DescribeHostRequest) (*host.Host, error) {
	//复用查询主机详情sql
	b := sqlbuilder.NewBuilder(QueryHostSql)
	b.Where("r.id=?", req.Id)
	stmt, args := b.Build()
	i.l.Debugf("describe sql:%v,args:%v", stmt, args)
	//数据库交互
	ins, err := i.describe(ctx, stmt, args)
	if err != nil {
		return nil, err
	}
	return ins, nil
}

//删除主机
func (i *HostServiceImpl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (*host.Host, error) {
	var err error
	//先获取要删除的host主机详情复用查询主机详情接口
	ins, err := i.DescribeHost(ctx, host.NewDescribeHostRequestWithId(req.Id))
	if err != nil {
		return nil, err
	}
	//执行删除host
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("start tx error:%s", err)
	}
	//通过defer处理事物逻辑
	//1.无错误则Commit事务
	//2.有报错则Rollback事务
	defer func() {
		//注意这里err判断的是ExecContext执行语句是否发生错误
		if err != nil {
			if err = tx.Rollback(); err != nil {
				i.l.Errorf("sql rollback error,%s", err)
			}
		} else {
			if err = tx.Commit(); err != nil {
				i.l.Errorf("sql commit error,%s", err)
			}
		}
	}()
	//删除resouce数据
	var (
		rstmt *sql.Stmt
		dstmt *sql.Stmt
	)
	rstmt, err = i.db.PrepareContext(ctx, DeleteResourceSql)
	if err != nil {
		return nil, err
	}
	defer rstmt.Close()
	if _, err = rstmt.ExecContext(ctx, req.Id); err != nil {
		return nil, err
	}
	//删除describe数据
	dstmt, err = i.db.PrepareContext(ctx, DeleteHostSql)
	if err != nil {
		return nil, err
	}
	defer dstmt.Close()
	if _, err = dstmt.ExecContext(ctx, req.Id); err != nil {
		return nil, err
	}
	return ins, nil
}

//注意这里指的空字段有2种含义1.没有设置值 2.设置了对应类型的零值
//在postman中如果想修改host对象中的某个字段，就要使用patch方法 补丁操作，传入非空的指定字段即可
//注意:如果请求对象存在字段值为空的情况(类型的零值如字符串为"" 数值为0)，服务端是不会更新此字段值的，这个时候就要使用put方法来进行全量更新
//patch中判断字段为空就不更新字段有2种情况：1.请求对象字段值设置成了空 2.请求对象中不包含的字段
//所以patch方法无法用于更新要将字段值置为空的请求，这个时候就要采用put来进行全量更新

//更新主机
func (i *HostServiceImpl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (*host.Host, error) {
	i.l.Error(req.Name)
	//先从数据库中获取要更新的对象
	ins, err := i.DescribeHost(ctx, host.NewDescribeHostRequestWithId(req.Id))
	if err != nil {
		return nil, err
	}
	//判断http方法
	switch req.HTTP_METHOD {
	case host.HTTP_METHOD_PUT:
		// ====注意dao层的sql已经指定允许的字段更新====
		//req数据对象与数据库中的ins对象操作,全量更新,req请求中未设置的字段使用其零值来覆盖ins对象字段
		if err := ins.Put(req.Host); err != nil {
			return nil, err
		}
	case host.HTTP_METHOD_PATCH:
		//req数据对象与数据库中的ins对象操作,局部更新，req带字段的值就覆盖对应字段
		i.l.Error(req.Name)
		if err := ins.Patch(req.Host); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("only support put/patch method!")
	}
	//检查更新后的数据是否合法 以防请求传入空值(对应字段的零值如int 0)
	i.l.Error(ins.Name)
	if err := ins.Validate(); err != nil {
		i.l.Error(err)
		return nil, err
	}

	//更新数据库
	if err := i.update(ctx, ins); err != nil {
		return nil, err
	}
	return ins, nil
}
