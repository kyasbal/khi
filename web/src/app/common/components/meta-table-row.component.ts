import { Component, Input } from '@angular/core';

@Component({
  selector: 'khi-meta-table-row',
  templateUrl: './meta-table-row.component.html',
  styleUrls: ['./meta-table-row.component.sass'],
})
export class MetaTableRowComponent {
  @Input()
  key = '';

  @Input()
  value = '';

  @Input()
  icon = '';
}
