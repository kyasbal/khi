import { TestBed } from '@angular/core/testing';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { ActivatedRouteSnapshot, RouterModule, UrlTree } from '@angular/router';
import { SessionChildGuard, SessionHostGuard } from './app.route.guard';
import { lastValueFrom } from 'rxjs';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import {
  WINDOW_CONNECTION_PROVIDER,
  WindowConnectionProvider,
  WindowConnectorService,
} from './services/frame-connection/window-connector.service';
import { InMemoryWindowConnectionProvider } from './services/frame-connection/window-connection-provider.service';

function createActivateRouteSnapshotWithSessionId(
  sessionId: string,
): ActivatedRouteSnapshot {
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  const snapshot = new (ActivatedRouteSnapshot as any)([], { sessionId });
  return snapshot;
}

describe('SessionHostGuard', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [RouterModule, NoopAnimationsModule, MatSnackBarModule],
      providers: [
        WindowConnectorService,
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
      ],
    });
  });

  it('should redirect to / when non-number session ID was given', async () => {
    const route = createActivateRouteSnapshotWithSessionId('foo');

    const guardResult = TestBed.runInInjectionContext(() =>
      SessionHostGuard(route),
    );

    const redirectTo = (await lastValueFrom(guardResult)) as UrlTree;
    expect(redirectTo.toString()).toBe('/session/0');
  });

  it('should accept routing when no other window have same session', async () => {
    const route = createActivateRouteSnapshotWithSessionId('0');

    const guardResult = TestBed.runInInjectionContext(() =>
      SessionHostGuard(route),
    );

    const redirectTo = (await lastValueFrom(guardResult)) as boolean;
    expect(redirectTo).toBe(true);
  });

  it('should redirect with incrementing sessionId when the given session ID is used', async () => {
    const route = createActivateRouteSnapshotWithSessionId('0');
    const connectionProvider = TestBed.inject<WindowConnectionProvider>(
      WINDOW_CONNECTION_PROVIDER,
    );
    const connector = new WindowConnectorService(connectionProvider);
    await connector.createSession(0);

    const guardResult = TestBed.runInInjectionContext(() =>
      SessionHostGuard(route),
    );

    const redirectTo = (await lastValueFrom(guardResult)) as UrlTree;
    expect(redirectTo.toString()).toBe('/session/1');
  });
});

describe('SessionChildGuard', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [MatSnackBarModule, RouterModule, NoopAnimationsModule],
      providers: [
        WindowConnectorService,
        {
          provide: WINDOW_CONNECTION_PROVIDER,
          useValue: new InMemoryWindowConnectionProvider(),
        },
      ],
    });
  });

  it('should redirect to / when non-number session ID was given', async () => {
    const route = createActivateRouteSnapshotWithSessionId('foo');

    const guardResult = TestBed.runInInjectionContext(() =>
      SessionHostGuard(route),
    );

    const redirectTo = (await lastValueFrom(guardResult)) as UrlTree;
    expect(redirectTo.toString()).toBe('/session/0');
  });

  it('should accept when the main window with same session ID is existing', async () => {
    const route = createActivateRouteSnapshotWithSessionId('11');
    const connectionProvider = TestBed.inject<WindowConnectionProvider>(
      WINDOW_CONNECTION_PROVIDER,
    );
    const connector = new WindowConnectorService(connectionProvider);
    await connector.createSession(11);

    const guardResult = TestBed.runInInjectionContext(() =>
      SessionChildGuard('Diagram')(route),
    );

    const redirectTo = (await lastValueFrom(guardResult)) as boolean;
    expect(redirectTo).toBe(true);
  });

  it('should redirect to /session/0 when the session ID main window is missing', async () => {
    const route = createActivateRouteSnapshotWithSessionId('11');

    const guardResult = TestBed.runInInjectionContext(() =>
      SessionChildGuard('Diagram')(route),
    );

    const redirectTo = (await lastValueFrom(guardResult)) as UrlTree;
    expect(redirectTo.toString()).toBe('/session/0');
  });
});
