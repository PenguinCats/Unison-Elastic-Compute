# UEC Redis 数据结构

## container
+ uec:container:{container_ID}:busy
  + type: string
  + 存在则表示繁忙
+ uec:container:{container_ID}:slave_ID
  + type: string
+ uec:container:{container_ID}:profile.dict
  + type: hash
  + hash fields
    + ext_container_id 
    + image_name
    + core_request
    + memory_request
    + storage_request
+ uec:container:{container_ID}:profile.exposed_tcp_ports
  + type: list
+ uec:container:{container_ID}:profile.exposed_tcp_mapping_ports
  + type: list
+ uec:container:{container_ID}:profile.exposed_udp_ports
  + type: list
+ uec:container:{container_ID}:profile.exposed_udp_mapping_ports
  + type: list
+ uec:container:{container_ID}:status
  + type: hash
  + hash fields
    + status
    + cpu_percent
    + memory_percent
    + memory_size
    + storage_size
    + storage_percent