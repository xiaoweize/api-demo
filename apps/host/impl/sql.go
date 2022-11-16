package impl

const (
	InsertResourceSQL = `
	INSERT INTO resource (
		id,
		vendor,
		region,
		create_at,
		expire_at,
		type,
		name,
		description,
		status,
		update_at,
		sync_at,
		accout,
		public_ip,
		private_ip
	)
	VALUES
		(?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`
	InsertDescribeSQL = `
	INSERT INTO host ( resource_id, cpu, memory, gpu_amount, gpu_spec, os_type, os_name, serial_number )
	VALUES
		( ?,?,?,?,?,?,?,? );
	`
	//注意下面sql作为sqlbuilder的基础sql后面不要加分号;
	//注意表名不要大写
	QueryHostSql = `
	SELECT
	r.*,h.cpu,h.memory,h.gpu_spec,h.gpu_amount,h.os_type,h.os_name,h.serial_number
	FROM
		resource AS r	
		LEFT JOIN host AS h ON r.id = h.resource_id
	`
	DeleteResourceSql = `
	DELETE from resource where id = ?
	`
	DeleteHostSql = `
	DELETE from host where resource_id = ?
	`

	//update sql
	updateResourceSQL = `UPDATE resource SET vendor=?,region=?,expire_at=?,name=?,description=? WHERE id = ?`
	updateHostSQL     = `UPDATE host SET cpu=?,memory=? WHERE resource_id = ?`
)
