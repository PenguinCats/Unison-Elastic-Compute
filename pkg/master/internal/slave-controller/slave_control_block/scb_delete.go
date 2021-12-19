package slave_control_block

func (scb *SlaveControlBlock) Delete() {
	_ = scb.RedisDAO.SlaveDelete(scb.uuid)
}
