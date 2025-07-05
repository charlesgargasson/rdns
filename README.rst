####
rDNS
####

| Golang PTR scanner

|

.. code-block:: bash

    $ ./bin/rdns --dns 192.168.1.1:53 --cidr 192.168.0.0/16
    Workers: 64, Timeout: 2s, DNS: 192.168.1.1:53
    Press Enter to stop scanning...

    CIDR 192.168.0.0/16 (65536 IPs)
    192.168.1.1 1.1.168.192.in-addr.arpa.	0	IN	PTR	livebox.home.
    192.168.1.18 18.1.168.192.in-addr.arpa.	0	IN	PTR	host.home.
    192.168.1.19 19.1.168.192.in-addr.arpa.	0	IN	PTR	esp32-368a58.home.
    192.168.1.20 20.1.168.192.in-addr.arpa.	0	IN	PTR	raspberrypi.home.
    [...]

    Scan completed in 15.380416328s 

.. code-block:: bash

    /opt/app # ./rdns --dns 10.96.0.10:53 --cidr k8s --workers 64
    Workers: 64, Timeout: 2s, DNS: 10.96.0.10:53
    Press Enter to stop scanning...

    CIDR 10.96.0.0/12 (1048576 IPs)
    10.96.0.1 1.0.96.10.in-addr.arpa.	30	IN	PTR	kubernetes.default.svc.cluster.local.
    10.96.0.10 10.0.96.10.in-addr.arpa.	30	IN	PTR	kube-dns.kube-system.svc.cluster.local.
    10.101.161.88 88.161.101.10.in-addr.arpa.	30	IN	PTR	grafana.adm.svc.cluster.local.
    10.101.234.203 203.234.101.10.in-addr.arpa.	30	IN	PTR	webapp.dev.svc.cluster.local.
    10.107.116.209 209.116.107.10.in-addr.arpa.	30	IN	PTR	registry.kube-system.svc.cluster.local.
    10.110.91.155 155.91.110.10.in-addr.arpa.	30	IN	PTR	ingress-nginx-controller-admission.ingress-nginx.svc.cluster.local.
    10.110.125.17 17.125.110.10.in-addr.arpa.	30	IN	PTR	ingress-nginx-controller.ingress-nginx.svc.cluster.local.

    CIDR 10.100.0.0/16 (65536 IPs)

    CIDR 10.0.0.0/16 (65536 IPs)

    CIDR 172.20.0.0/16 (65536 IPs)

    Scan completed in 37.938820541s 

|
