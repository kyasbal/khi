import { Observable, Subject, firstValueFrom } from 'rxjs';
import { InterframeDatasource } from './inter-frame-datasource.service';

describe('InterframeDatasource', () => {
  class TestingInterframeDataSource extends InterframeDatasource<string> {
    override enable(): void {
      return;
    }
    override disable(): void {
      return;
    }

    constructor(dataSource: Observable<string>) {
      super();
      dataSource.subscribe(this.rawUpdateRequest$);
    }
  }

  it('should emit the data when observer is registered', async () => {
    const dataSourceSubject = new Subject<string>();
    const datasource = new TestingInterframeDataSource(dataSourceSubject);

    dataSourceSubject.next('foo');
    dataSourceSubject.next('bar');

    expect(await firstValueFrom(datasource.data$)).toBe('bar');
  });

  it('should not emit the data when bound$ is false', async () => {
    const dataSourceSubject = new Subject<string>();
    const datasource = new TestingInterframeDataSource(dataSourceSubject);

    dataSourceSubject.next('foo');
    datasource.bound$.next(false);
    dataSourceSubject.next('bar');

    expect(await firstValueFrom(datasource.data$)).toBe('foo');
  });
});
