terraform {
    required_providers {
        kubernetes = {
            source = "hashicorp/kubernetes"
            version = "2.25.2"
        }
        helm = {
            source = "hashicorp/helm"
            version = "2.12.1"
        }
    }
}

provider "helm" {
    kubernetes {
        config_path = "~/.kube/config"
    }
}

provider "kubernetes" {
    config_path = "~/.kube/config"
    config_context = "minikube"
}

resource "kubernetes_namespace" "tools_namespace" {
  metadata {
    name = "monitoring"
  }
}

resource "helm_release" "kube-prometheus-stack" {
  name       = "kube-prometheus-stack"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"
  version    = "56.7.0"
  
  namespace  = "monitoring"

  values = [
    "${file("${path.module}/kube-prometheus-stack.yaml")}"
  ]
}


resource "helm_release" "grafana" {
  name       = "grafana"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "grafana"
  version    = "7.3.2"

  namespace  = "monitoring"

  values = [
    "${file("${path.module}/grafana.yaml")}"
  ]
  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}


resource "helm_release" "prometheus-openweathermap-exporter" {
  name      = "owm"
  chart = "${path.module}/../charts/prometheus-openweathermap-exporter"
  version = "1.1.0"
  namespace = "monitoring"
  values = [
    "${file("${path.module}/prometheus-openweathermap-exporter.yaml")}"
  ]
  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}
