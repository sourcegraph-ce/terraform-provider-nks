resource "nks_solution" "jenkins" {
  org_id     = "${data.nks_keysets.keyset_default.org_id}"
  cluster_id = "${nks_cluster.terraform-cluster.id}"
  solution   = "jenkins"
  config     = "${file("solutions/jenkins.json")}"
}
