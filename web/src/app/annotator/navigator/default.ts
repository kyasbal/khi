import { AnnotationDecider, Annotator, DECISION_HIDDEN } from '../annotator';
import { CommonFieldAnnotatorComponent } from '../common-field-annotator.component';
import {
  CommonReferenceListAnnotatorDecision,
  CommonReferenceListComponent,
  ReferenceViewModel,
} from '../common-reference-list.component';
import { NavigatorAnnotatorResolver } from './resolver';
import { of } from 'rxjs';
import {
  K8sNodeResource,
  K8sPodBindingResource,
  K8sPodResource,
} from 'src/app/models/k8s/k8s-types';
import { componentDocuments } from './doc';
import { ResourceRevision } from 'src/app/store/revision';
import { TimelineEntry, TimelineLayer } from 'src/app/store/timeline';

type referenceListMatcher = (rev: ResourceRevision) => ReferenceViewModel[];

function referenceListFromRevisions(
  ...matchers: referenceListMatcher[]
): AnnotationDecider<TimelineEntry> {
  return (t) => {
    if (!t) return DECISION_HIDDEN;
    const contained = new Set();
    const result: ReferenceViewModel[] = [];
    for (const matcher of matchers) {
      for (const revision of t.revisions) {
        const references = matcher(revision);
        if (!Array.isArray(references)) continue;
        for (const ref of references) {
          if (!contained.has(ref.url)) {
            contained.add(ref.url);
            result.push(ref);
          }
        }
      }
    }
    if (result.length === 0) return DECISION_HIDDEN;
    return {
      inputs: {
        references: of(result),
      },
    } as CommonReferenceListAnnotatorDecision;
  };
}

const COMPONENT_MATCHER: referenceListMatcher = (rev: ResourceRevision) =>
  componentDocuments[
    rev.parsedManifest?.metadata?.annotations?.[
      'components.gke.io/component-name'
    ] ?? ''
  ];

function nodeNameFromTimelineAnnotationDecider(): AnnotationDecider<TimelineEntry> {
  const icon = 'dns';
  const label = 'Node';
  return (tl) => {
    if (!tl) return DECISION_HIDDEN;
    if (
      tl.layer !== TimelineLayer.Name ||
      tl.getNameOfLayer(TimelineLayer.Kind) !== 'pod'
    )
      return DECISION_HIDDEN;
    let nodeNames = tl.revisions.map(
      (rev) => (rev.parsedManifest as K8sPodResource)?.spec?.nodeName,
    );
    const bindings = tl.children.filter(
      (tl) => tl.getNameOfLayer(TimelineLayer.Subresource) === 'binding',
    );
    if (bindings.length > 0) {
      const binding = bindings[0];
      nodeNames.push(
        ...binding.revisions.map(
          (rev) => (rev.parsedManifest as K8sPodBindingResource)?.target?.name,
        ),
      );
    }
    nodeNames = nodeNames.filter((n) => !!n);
    if (nodeNames.length === 0)
      return {
        inputs: {
          icon,
          label,
          value: of(['Unknown']),
        },
      };
    return {
      inputs: {
        icon,
        label,
        value: of([[...new Set(nodeNames)].join(',')]),
      },
    };
  };
}

function evedashboardForNodeAnnotationDecider(): AnnotationDecider<TimelineEntry> {
  return (tl) => {
    if (!tl) return DECISION_HIDDEN;
    if (
      tl.layer !== TimelineLayer.Name ||
      tl.getNameOfLayer(TimelineLayer.Kind) !== 'node'
    ) {
      return DECISION_HIDDEN;
    }
    const version = tl.revisions
      .map(
        (rev) =>
          (rev.parsedManifest as K8sNodeResource)?.status?.nodeInfo
            ?.kubeletVersion,
      )
      .filter((n) => !!n);
    if (version.length === 0) {
      return DECISION_HIDDEN;
    } else {
      return {
        inputs: {
          header: 'Eve',
          references: of(
            [...new Set(version)].map((v) => ({
              url: `https://evedashboard.corp.google.com/versions/${v?.substring(1)}`,
              displayText: `${v}`,
            })),
          ),
        },
      } as CommonReferenceListAnnotatorDecision;
    }
  };
}

function evedashboardForComponentAnnotationDecider(): AnnotationDecider<TimelineEntry> {
  return (tl) => {
    if (!tl) return DECISION_HIDDEN;
    if (
      tl.layer !== TimelineLayer.Name ||
      tl.getNameOfLayer(TimelineLayer.Kind) !== 'pod'
    ) {
      return DECISION_HIDDEN;
    }
    const componentNames = tl.revisions
      .map(
        (rev) =>
          (rev.parsedManifest as K8sPodResource)?.metadata?.annotations?.[
            'components.gke.io/component-name'
          ],
      )
      .filter((n) => !!n);
    if (componentNames.length === 0) return DECISION_HIDDEN;
    const componentName = componentNames[0];
    const componentVersions = tl.revisions
      .map(
        (rev) =>
          (rev.parsedManifest as K8sPodResource)?.metadata?.annotations?.[
            'components.gke.io/component-version'
          ],
      )
      .filter((n) => !!n);
    if (componentVersions.length === 0) {
      return DECISION_HIDDEN;
    } else {
      return {
        inputs: {
          header: 'Eve',
          references: of(
            [...new Set(componentVersions)].map((v) => ({
              url: `https://evedashboard.corp.google.com/components/${componentName}/versions/${v}`,
              displayText: `${v}`,
            })),
          ),
        },
      } as CommonReferenceListAnnotatorDecision;
    }
  };
}

export function getDefaultNavigatorAnnotatorResolver(): NavigatorAnnotatorResolver {
  return new NavigatorAnnotatorResolver([
    new Annotator(
      CommonFieldAnnotatorComponent,
      nodeNameFromTimelineAnnotationDecider(),
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForTimelineOfRevisions(
        '',
        'Component',
        (rev) =>
          rev.parsedManifest?.metadata?.annotations?.[
            'components.gke.io/component-name'
          ],
      ),
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForTimelineOfRevisions(
        '',
        'Version',
        (rev) =>
          rev.parsedManifest?.metadata?.annotations?.[
            'components.gke.io/component-version'
          ],
      ),
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForTimelineOfRevisions(
        '',
        'Instance type',
        (rev) =>
          rev.parsedManifest?.metadata?.labels?.[
            'node.kubernetes.io/instance-type'
          ],
      ),
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForTimelineOfRevisions(
        '',
        'Zone',
        (rev) =>
          rev.parsedManifest?.metadata?.labels?.['topology.kubernetes.io/zone'],
      ),
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForTimelineOfRevisions(
        '',
        'OS',
        (rev) =>
          rev.parsedManifest?.metadata?.labels?.[
            'cloud.google.com/gke-os-distribution'
          ],
      ),
    ),
    new Annotator(
      CommonReferenceListComponent,
      referenceListFromRevisions(COMPONENT_MATCHER),
    ),
    new Annotator(
      CommonReferenceListComponent,
      evedashboardForNodeAnnotationDecider(),
    ),
    new Annotator(
      CommonReferenceListComponent,
      evedashboardForComponentAnnotationDecider(),
    ),
  ]);
}
