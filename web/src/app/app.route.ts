import { Routes } from '@angular/router';
import { AppComponent } from './pages/main/main.component';
import { GraphComponent } from './pages/graph/graph.component';
import {
  DiffPageDeactivateGuard,
  DiffPageGuard,
  FrontendAnalyticsServiceGuard,
  GraphPageDeactiveGuard,
  GraphPageGuard,
  SessionChildGuard,
  SessionDeactivateGuard,
  SessionHostGuard,
} from './app.route.guard';
import { DiffComponent } from './pages/diff/diff.component';
import { KHIAnalyticsPageType } from './services/analytics/types';

export const KHIRoutes: Routes = [
  { path: '', redirectTo: 'session/0', pathMatch: 'full' },
  {
    path: 'session/:sessionId',
    component: AppComponent,
    title: 'KHI - Main view',
    canActivate: [
      SessionHostGuard,
      FrontendAnalyticsServiceGuard(KHIAnalyticsPageType.Main),
    ],
    canDeactivate: [SessionDeactivateGuard],
  },
  {
    path: 'session/:sessionId/graph',
    component: GraphComponent,
    title: 'KHI - Graph view',
    canActivate: [
      SessionChildGuard('Diagram'),
      GraphPageGuard,
      FrontendAnalyticsServiceGuard(KHIAnalyticsPageType.GraphView),
    ],
    canDeactivate: [SessionDeactivateGuard, GraphPageDeactiveGuard],
  },
  {
    path: 'session/:sessionId/diff/:kind/:namespace/:resourceName/:subresource',
    component: DiffComponent,
    title: 'KHI - Diff view',
    canActivate: [
      SessionChildGuard('Diff'),
      DiffPageGuard,
      FrontendAnalyticsServiceGuard(KHIAnalyticsPageType.DiffView),
    ],
    canDeactivate: [SessionDeactivateGuard, DiffPageDeactivateGuard],
  },
  { path: '**', redirectTo: '/' },
];
