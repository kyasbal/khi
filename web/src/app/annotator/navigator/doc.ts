import { ReferenceViewModel } from '../common-reference-list.component';

// TODO: This is experimental feature. This list would be better to be managed in Piper.

export const componentDocuments: { [component: string]: ReferenceViewModel[] } =
  {
    'event-exporter': [
      {
        url: 'http://go/gke-event-exporter#event-exporter-in-a-gke-cluster',
        displayText: 'Playbook',
      },
    ],
    fluentbit: [
      {
        url: 'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/container-engine/gke-logging-troubleshooting-tree-playbooks/gke-logging-fluentbit-issues.md?cl=head',
        displayText: 'Playbook',
      },
      {
        url: 'https://fluentbit.io/',
        displayText: 'External document(fluentbit)',
      },
    ],
    'gke-metrics-agent': [
      {
        url: 'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/container-engine/metrics-agent.md?cl=head',
        displayText: 'Playbook',
      },
    ],
    'konnectivitynetworkproxy-combined': [
      {
        displayText: 'Playbook',
        url: 'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/container-engine/konnectivity.md?cl=head',
      },
    ],
    kubedns: [
      {
        url: 'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/kubernetes-networking/kube-dns-debug-tree.md?cl=head',
        displayText: 'Playbook(Infra)',
      },
      {
        url: 'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/kubernetes-networking/kubedns.md?cl=head',
        displayText: 'Playbook(Network)',
      },
    ],
    'l7-lb-controller-combined': [
      {
        url: 'https://playbooks.corp.google.com/cloud-kubernetes/components/l7-lb-controller.md?cl=head',
        displayText: 'Internal doc',
      },
    ],
    'metrics-server': [
      {
        url: 'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/container-engine/metrics-server.md?cl=head',
        displayText: 'Playbook',
      },
    ],
    pdcsi: [
      {
        url: 'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/container-engine/gke-storage.md?cl=head',
        displayText: 'Playbook',
      },
    ],
    'managed-prometheus': [
      {
        url: 'https://g3doc.corp.google.com/company/gfw/support/cloud/playbooks/monitoring/gmp-troubleshooting.md?cl=head',
        displayText: 'Playbook',
      },
    ],
  };
