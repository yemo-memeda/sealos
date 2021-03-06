package install

import (
	"fmt"
	"sync"
)

//BuildClean is
func BuildClean(beforeNodes []string) {
	i := &SealosInstaller{}
	var hosts []string
	if len(beforeNodes) == 0 {
		// 所有master节点
		masters := ParseIPs(MasterIPs)
		// 所有node节点
		nodes := ParseIPs(NodeIPs)
		hosts = append(masters, nodes...)
	}else {
		hosts = beforeNodes
	}
	i.Hosts = hosts
	i.CheckValid()
	i.Clean()
}

//CleanCluster is
func (s *SealosInstaller) Clean() {
	var wg sync.WaitGroup
	for _, host := range s.Hosts {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()
			clean(node)
		}(host)
	}
	wg.Wait()
}

func clean(host string) {
	cmd := "kubeadm reset -f && modprobe -r ipip  && lsmod"
	Cmd(host, cmd)
	cmd = "rm -rf ~/.kube/ && rm -rf /etc/kubernetes/"
	Cmd(host, cmd)
	cmd = "rm -rf /etc/systemd/system/kubelet.service.d && rm -rf /etc/systemd/system/kubelet.service"
	Cmd(host, cmd)
	cmd = "rm -rf /usr/bin/kube* && rm -rf /usr/bin/crictl"
	Cmd(host, cmd)
	cmd = "rm -rf /etc/cni && rm -rf /opt/cni"
	Cmd(host, cmd)
	cmd = "rm -rf /var/lib/etcd && rm -rf /var/etcd"
	Cmd(host, cmd)
	cmd = fmt.Sprintf("sed -i \"/%s/d\" /etc/hosts ", ApiServer)
	Cmd(host, cmd)
	cmd = fmt.Sprint("rm -rf ~/kube")
	Cmd(host, cmd)
}
