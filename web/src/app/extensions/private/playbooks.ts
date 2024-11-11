/**
 * Copyright 2024 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

export interface PlaybookBindingWithComponentNameAnnotationData {
  componentName: string;
  linkText: string;
  linkUrl: string;
}

export const PLAYBOOK_DATA: PlaybookBindingWithComponentNameAnnotationData[] = [
  {
    componentName: 'event-exporter',
    linkText: '"event-exporter" playbook',
    linkUrl: 'http://go/gke-event-exporter#event-exporter-in-a-gke-cluster',
  },
  {
    componentName: 'fluentbit',
    linkText: '"fluentbit" playbook',
    linkUrl:
      'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/container-engine/gke-logging-troubleshooting-tree-playbooks/gke-logging-fluentbit-issues.md?cl=head',
  },
  {
    componentName: 'gke-metrics-agent',
    linkText: '"gke-metrics-agent" playbook',
    linkUrl:
      'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/container-engine/metrics-agent.md?cl=head',
  },
  {
    componentName: 'konnectivitynetworkproxy-combined',
    linkText: '"konnectivity-agent" playbook',
    linkUrl:
      'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/container-engine/konnectivity.md?cl=head',
  },
  {
    componentName: 'kubedns',
    linkText: '"kubedns" playbook for infra',
    linkUrl:
      'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/kubernetes-networking/kube-dns-debug-tree.md?cl=head',
  },
  {
    componentName: 'kubedns',
    linkText: '"kubedns" playbook for network',
    linkUrl:
      'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/kubernetes-networking/kubedns.md?cl=head',
  },
  {
    componentName: 'l7-lb-controller-combined',
    linkText: '"l7-lb-controller" playbook for SREs',
    linkUrl:
      'https://playbooks.corp.google.com/cloud-kubernetes/components/l7-lb-controller.md?cl=head',
  },
  {
    componentName: 'metrics-server',
    linkText: '"metrics-server" playbook',
    linkUrl:
      'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/container-engine/metrics-server.md?cl=head',
  },
  {
    componentName: 'pdcsi',
    linkText: '"gke-storage" playbook',
    linkUrl:
      'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/container-engine/gke-storage.md?cl=head',
  },
  {
    componentName: 'managed-prometheus',
    linkText: '"GMP" playbook',
    linkUrl:
      'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/monitoring/gmp-troubleshooting.md?cl=head',
  },
];
