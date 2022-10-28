package impl

import (
	"context"
	"fmt"

	"github.com/xiaoweize/api-demo/apps/host"
)

//go对象与数据库之间的转换

//将host.Host对象保存到mysql数据库
func (i *HostServiceImpl) save(ctx context.Context, ins *host.Host) error {
	// 把数据入库到 resource表和host表
	// 一次需要往2个表录入数据, 我们需要2个操作 要么都成功，要么都失败,  事务的逻辑
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start tx error %s", err)
	}
	//通过defer处理事物逻辑 
	//1.无错误则Commit事物
	//2.有报错则Rollback事物
	defer func() {
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
	//插入resource数据
	rstmt, err := tx.Prepare(InsertResourceSQL)
	if err != nil {
		return err
	}
	//每一个statment 都会维持在mysql的内存中，如果不关闭，会导致数据库资源一直被占用，所以一定要释放
	defer rstmt.Close()
	_, err = rstmt.Exec(
		ins.Id, ins.Vendor, ins.Region, ins.CreateAt, ins.ExpireAt, ins.Type,
		ins.Name, ins.Description, ins.Status, ins.UpdateAt, ins.SyncAt, ins.Account, ins.PublicIP,
		ins.PrivateIP,
	)
	if err != nil {
		return err
	}

	//插入Describe数据
	dstmt, err := tx.Prepare(InsertDescribeSQL)
	if err != nil {
		return err
	}
	//每一个statment 都会维持在mysql的内存中，如果不关闭，会导致数据库资源一直被占用，所以一定要释放
	defer dstmt.Close()
	_, err = dstmt.Exec(
		ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec,
		ins.OSType, ins.OSName, ins.SerialNumber,
	)
	if err != nil {
		return err
	}
	return nil
}
