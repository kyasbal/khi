import { ParsePrincipalPipe, PrincipalType } from './diff-view-pipes';

describe('ParsePrincipalPipe', () => {
  it('user account', () => {
    const ppp = new ParsePrincipalPipe();
    const principal = ppp.transform('foo@bar.com');
    expect(principal.type).toBe(PrincipalType.User);
    expect(principal.short).toBe('foo@bar.com');
    expect(principal.full).toBe('foo@bar.com');
  });
  it('system account', () => {
    const ppp = new ParsePrincipalPipe();
    const principal = ppp.transform('system:garbage-collector');
    expect(principal.type).toBe(PrincipalType.System);
    expect(principal.short).toBe('garbage-collector');
    expect(principal.full).toBe('system:garbage-collector');
  });
  it('node account', () => {
    const ppp = new ParsePrincipalPipe();
    const principal = ppp.transform('system:node:node-foo-bar');
    expect(principal.type).toBe(PrincipalType.Node);
    expect(principal.short).toBe('node-foo-bar');
    expect(principal.full).toBe('system:node:node-foo-bar');
  });
  it('service account', () => {
    const ppp = new ParsePrincipalPipe();
    const principal = ppp.transform('system:serviceaccount:argocd');
    expect(principal.type).toBe(PrincipalType.ServiceAccount);
    expect(principal.short).toBe('argocd');
    expect(principal.full).toBe('system:serviceaccount:argocd');
  });
});
