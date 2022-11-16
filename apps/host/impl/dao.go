package impl

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/xiaoweize/api-demo/apps/host"
)

//go对象与数据库之间的转换

//将host.Host对象保存到mysql数据库
//使用事务，保持一致性
func (i *HostServiceImpl) save(ctx context.Context, ins *host.Host) error {
	var err error
	// 把数据入库到 resource表和host表
	// 一次需要往2个表录入数据, 我们需要2个操作要么都成功，要么都失败,事务的逻辑
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start tx error %s", err)
	}
	//通过defer处理事物逻辑
	//1.无错误则Commit事务
	//2.有报错则Rollback事务
	defer func() {
		if err != nil {
			//注意这里err判断的是PrepareContext/ExecContext执行语句是否发生错误
			if err = tx.Rollback(); err != nil {
				i.l.Errorf("sql rollback error,%s", err)
			}
		} else {
			if err = tx.Commit(); err != nil {
				i.l.Errorf("sql commit error,%s", err)
			}
		}
	}()
	//插入resource数据
	var (
		rstmt *sql.Stmt
		dstmt *sql.Stmt
	)
	rstmt, err = tx.PrepareContext(ctx, InsertResourceSQL)
	if err != nil {
		return err
	}
	//每一个statment 都会维持在mysql的内存中，如果不关闭，会导致数据库资源一直被占用，所以一定要释放
	defer rstmt.Close()
	_, err = rstmt.ExecContext(ctx,
		ins.Id, ins.Vendor, ins.Region, ins.CreateAt, ins.ExpireAt, ins.Type,
		ins.Name, ins.Description, ins.Status, ins.UpdateAt, ins.SyncAt, ins.Account, ins.PublicIP,
		ins.PrivateIP,
	)
	if err != nil {
		return err
	}

	//插入Describe数据
	dstmt, err = tx.PrepareContext(ctx, InsertDescribeSQL)
	if err != nil {
		return err
	}
	//每一个statment 都会维持在mysql的内存中，如果不关闭，会导致数据库资源一直被占用，所以一定要释放
	defer dstmt.Close()
	_, err = dstmt.ExecContext(ctx,
		ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec,
		ins.OSType, ins.OSName, ins.SerialNumber,
	)
	if err != nil {
		return err
	}
	return nil
}

func (i *HostServiceImpl) describe(ctx context.Context, query string, args []any) (*host.Host, error) {
	stmt, err := i.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	ins := host.NewHost()
	if err := stmt.QueryRowContext(ctx, args...).Scan(&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt, &ins.ExpireAt,
		&ins.Type, &ins.Name, &ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt,
		&ins.Account, &ins.PublicIP, &ins.PrivateIP,
		&ins.CPU, &ins.Memory, &ins.GPUSpec, &ins.GPUAmount, &ins.OSType, &ins.OSName, &ins.SerialNumber); err != nil {
		return nil, err
	}
	return ins, nil
}

func (i *HostServiceImpl) update(ctx context.Context, ins *host.Host) error {
	var err error
	// 把数据更新到 resource表和host表
	// 一次需要往2个表更新数据, 我们需要2个操作要么都成功，要么都失败,事务的逻辑
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start tx error %s", err)
	}
	//通过defer处理事物逻辑
	//1.无错误则Commit事务
	//2.有报错则Rollback事务
	defer func() {
		if err != nil {
			//注意这里err判断的是PrepareContext/ExecContext执行语句是否发生错误
			if err = tx.Rollback(); err != nil {
				i.l.Errorf("sql rollback error,%s", err)
			}
		} else {
			if err = tx.Commit(); err != nil {
				i.l.Errorf("sql commit error,%s", err)
			}
		}
	}()
	//更新resource数据
	var (
		rstmt, dstmt *sql.Stmt
	)
	rstmt, err = tx.PrepareContext(ctx, updateResourceSQL)
	if err != nil {
		return err
	}
	//每一个statment 都会维持在mysql的内存中，如果不关闭，会导致数据库资源一直被占用，所以一定要释放
	defer rstmt.Close()
	_, err = rstmt.ExecContext(ctx, ins.Vendor, ins.Region, ins.ExpireAt, ins.Name, ins.Description, ins.Id)
	if err != nil {
		return err
	}

	//更新Describe数据
	dstmt, err = tx.PrepareContext(ctx, updateHostSQL)
	if err != nil {
		return err
	}
	//每一个statment 都会维持在mysql的内存中，如果不关闭，会导致数据库资源一直被占用，所以一定要释放
	defer dstmt.Close()
	_, err = dstmt.ExecContext(ctx, ins.CPU, ins.Memory, ins.Id)
	if err != nil {
		return err
	}
	return nil
}
