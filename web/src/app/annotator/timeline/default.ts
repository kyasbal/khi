import { TimelineLayer } from 'src/app/store/timeline';
import { Annotator } from '../annotator';
import { CommonFieldAnnotatorComponent } from '../common-field-annotator.component';
import { TimelineAnnotatorResolver } from './resolver';

export function getDefaultTimelineAnnotatorResolver(): TimelineAnnotatorResolver {
  return new TimelineAnnotatorResolver([
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForTimelineEntry(
        'workspaces',
        'Kind',
        (tl) => tl.getNameOfLayer(TimelineLayer.Kind),
      ),
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForTimelineEntry(
        'folder',
        'Namespace',
        (tl) => tl.getNameOfLayer(TimelineLayer.Namespace),
      ),
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForTimelineEntry(
        'description',
        'Name',
        (tl) => tl.getNameOfLayer(TimelineLayer.Name),
      ),
    ),
    new Annotator(
      CommonFieldAnnotatorComponent,
      CommonFieldAnnotatorComponent.inputMapperForTimelineEntry(
        'page_info',
        'Subresource',
        (tl) => tl.getNameOfLayer(TimelineLayer.Subresource),
      ),
    ),
  ]);
}
