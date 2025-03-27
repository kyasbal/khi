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

import { ResourceTimeline, TimelineLayer } from 'src/app/store/timeline';
import {
  DisplayableTimelineNavigatorExtension,
  TimelineNavigatorExtension,
  TimelineNavigatorExtensionUtil,
} from '../extension-common/extension-types/timeline-navigator';
import { CommonReferenceListComponent } from 'src/app/annotator/common-reference-list.component';
import { TimelineFilterFacade } from 'src/app/store/timeline-filter';

/**
 * PlaybookBindingWithComponentNameAnnotation is a TimelineNavigatorExtension showing playbook link with matching `metadata.annotations."components.gke.io/component-name"` annotation.
 */
export class PlaybookBindingWithComponentNameAnnotation
  implements TimelineNavigatorExtension
{
  constructor(
    private readonly componentName: string,
    private readonly linkText: string,
    private readonly linkUrl: string,
  ) {}

  show(timeline: ResourceTimeline): boolean {
    return (
      timeline.layer === TimelineLayer.Name &&
      TimelineNavigatorExtensionUtil.anyOfManifestBodyFieldInRevisions(
        timeline,
        ['metadata', 'annotations', 'components.gke.io/component-name'],
        (value) => value === this.componentName,
      )
    );
  }
  getDisplayable(): DisplayableTimelineNavigatorExtension {
    return {
      component: CommonReferenceListComponent,
      inputs: {
        references: [
          {
            linkText: this.linkText,
            linkUrl: this.linkUrl,
          },
        ],
      },
    };
  }
}

/**
 * EveDashboardBindingForNode is a TimelineNavigatorExtension showing EveDashboard for the specific node version.
 */
export class EveDashboardBindingForNode implements TimelineNavigatorExtension {
  show(timeline: ResourceTimeline): boolean {
    return (
      timeline.layer === TimelineLayer.Name &&
      TimelineFilterFacade.isNodeOrNodeChildren(timeline)
    );
  }
  getDisplayable(
    timeline: ResourceTimeline,
  ): DisplayableTimelineNavigatorExtension {
    const kubeletVersions =
      TimelineNavigatorExtensionUtil.getSetOfManifestBodyFieldInRevisions(
        timeline,
        ['status', 'nodeInfo', 'kubeletVersion'],
      );
    return {
      component: CommonReferenceListComponent,
      inputs: {
        references: kubeletVersions
          .map((v) => v + '') // converting unknown[] to string[]
          .map((v) => ({
            linkText: `Eve Dashboard(Node ${v})`,
            linkUrl: `https://evedashboard.corp.google.com/versions/${v.substring(1)}`,
          })),
      },
    };
  }
}

export class EveDashboardBindingForComponent
  implements TimelineNavigatorExtension
{
  show(timeline: ResourceTimeline): boolean {
    return (
      timeline.layer === TimelineLayer.Name &&
      TimelineNavigatorExtensionUtil.anyOfManifestBodyFieldInRevisions(
        timeline,
        ['metadata', 'annotations', 'components.gke.io/component-name'],
        (v) => !!v,
      ) &&
      TimelineNavigatorExtensionUtil.anyOfManifestBodyFieldInRevisions(
        timeline,
        ['metadata', 'annotations', 'components.gke.io/component-version'],
        (v) => !!v,
      )
    );
  }
  getDisplayable(
    timeline: ResourceTimeline,
  ): DisplayableTimelineNavigatorExtension {
    const componentNames =
      TimelineNavigatorExtensionUtil.getSetOfManifestBodyFieldInRevisions(
        timeline,
        ['metadata', 'annotations', 'components.gke.io/component-name'],
      );
    const componentName = componentNames[0]; // this array should have at least 1 element because it should be filtered on show() method
    const componentVersions =
      TimelineNavigatorExtensionUtil.getSetOfManifestBodyFieldInRevisions(
        timeline,
        ['metadata', 'annotations', 'components.gke.io/component-version'],
      );
    return {
      component: CommonReferenceListComponent,
      inputs: {
        references: componentVersions
          .map((v) => v + '') // converting unknown[] to string[]
          .map((v) => ({
            linkText: `Eve Dashboard(Component ${componentName}, version ${v})`,
            linkUrl: `https://evedashboard.corp.google.com/components/${componentName}/versions/${v}`,
          })),
      },
    };
  }
}
