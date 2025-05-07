# 混合云场景下物理专线是通过OPSAPI创建的，通过Terraform创建物理专线时，provider需要使用运维端账号的AKSK
resource "vestack_direct_connect_connection" "foo"{
  direct_connect_connection_name="tf-test-connection"
  description="tf-test"
  port_id="dcp-xxxxx"
  # 运维端不持有资源，owner_account_id填入租户的account_id
  owner_account_id="1000000xxx"
  owner_project_name="default"
  line_operator="ChinaOther"
  port_type="10GBase"
  port_spec="10G"
  bandwidth=1000
  peer_location="XX路XX号XX楼XX机房"
  customer_name="tf-a"
  customer_contact_phone="12345678911"
  customer_contact_email="email@aaa.com"
}