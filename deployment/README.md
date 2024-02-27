# Deployment

This directory contains the Terraform code to deploy the exporter and related tools to
a Minikube cluster. It will install the following charts:

- Kube-Prometheus-Stack
- Grafana
- Prometheus-Openweathermap-Exporter

## Usage

Running the Prometheus-Openweathermap-Exporter requires an OWM API key. To provide it,
create a `prometheus-openweathermap-exporter.yaml` values file in this directory:

```bash
$ echo "apiKey: <redacted>" >> ./prometheus-openweathermap-exporter.yaml
...
```

Start a Minikube cluster, then apply the Terraform code:

```bash
$ minikube start
...
$ terraform apply
...
```

Prometheus and Grafana have their Ingress enabled. Using `minikube tunnel` and editing your `/etc/hosts` file,
you should be able to access them via browser

```bash
$ kubectl get ingress -n monitoring
NAME                               CLASS   HOSTS              ADDRESS        PORTS   AGE
grafana                            nginx   grafana.local      192.168.49.2   80      22s
kube-prometheus-stack-prometheus   nginx   prometheus.local   192.168.49.2   80      41s

# as root
$ echo "192.168.42.2  grafana.local" >> /etc/hosts
$ echo "192.168.42.2  prometheus.local" >> /etc/hosts
$ minikube tunnel
```

On Windows you may need to use `127.0.0.1` instead in your hosts file.
