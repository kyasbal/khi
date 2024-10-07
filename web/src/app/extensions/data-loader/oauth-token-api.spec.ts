import { Subject } from 'rxjs';
import { OAuthTokenAPI } from './oauth-token-api';

describe('OAuthTokenAPI', () => {
  class MockLocalStorage {
    constructor(private data: { [key: string]: string | null } = {}) {}

    public setItem(key: string, value: string | null) {
      this.data[key] = value;
    }

    public removeItem(key: string) {
      delete this.data[key];
    }

    public getItem(key: string): string | null {
      return this.data[key] ?? null;
    }
  }

  function createOAuthAPI(
    control: {
      currentAccessToken: string | null;
    } = { currentAccessToken: null },
  ): OAuthTokenAPI {
    const onRequestAccessTokenCalled: Subject<void> = new Subject();
    const oauthClientSpy = {
      requestAccessToken: () => {
        onRequestAccessTokenCalled.next();
      },
    };
    const oauthRefSpy = jasmine.createSpyObj('oauth2', ['initTokenClient']);
    oauthRefSpy.initTokenClient.and.returnValue(oauthClientSpy);

    const oauthAPI = new OAuthTokenAPI(
      oauthRefSpy,
      new MockLocalStorage() as unknown as Storage,
    );
    const tokenCallback =
      oauthRefSpy.initTokenClient.calls.first().args[0].callback;
    onRequestAccessTokenCalled.subscribe(() =>
      tokenCallback({ access_token: control.currentAccessToken }),
    );
    oauthAPI.clearTokenCache();
    return oauthAPI;
  }

  describe('#requestAccessToken', () => {
    it('should throw error if token was null', async () => {
      const oauthAPI = createOAuthAPI();
      await expectAsync(oauthAPI.requestAccessToken(false)).toBeRejected();
    });

    it('should return token if valid token was passed', async () => {
      const oauthAPI = createOAuthAPI({ currentAccessToken: 'TESTTOKEN' });
      await expectAsync(oauthAPI.requestAccessToken(false)).toBeResolvedTo(
        'TESTTOKEN',
      );
    });

    it('should use cached token if cached token is available', async () => {
      const currentTokenInfo = { currentAccessToken: 'TESTTOKEN' };
      const oauthAPI = createOAuthAPI(currentTokenInfo);
      await oauthAPI.requestAccessToken(false);
      currentTokenInfo.currentAccessToken = 'INVALIDTOKEN';
      await expectAsync(oauthAPI.requestAccessToken(false)).toBeResolvedTo(
        'TESTTOKEN',
      );
    });

    it("shouldn't use cached token if cached token is available but disabled", async () => {
      const currentTokenInfo = { currentAccessToken: 'TESTTOKEN' };
      const oauthAPI = createOAuthAPI(currentTokenInfo);
      await oauthAPI.requestAccessToken(true);
      currentTokenInfo.currentAccessToken = 'TESTTOKEN2';
      await expectAsync(oauthAPI.requestAccessToken(true)).toBeResolvedTo(
        'TESTTOKEN2',
      );
    });
  });

  describe('#requestAPIWithOAuth', () => {
    it('should call api with token', async () => {
      const currentTokenInfo = { currentAccessToken: 'TESTTOKEN' };
      const oauthAPI = createOAuthAPI(currentTokenInfo);
      const spyAPI = jasmine
        .createSpy('callAPI')
        .and.returnValue(Promise.resolve('SUCCESS'));
      await expectAsync(
        oauthAPI.requestAPIWithOAuth(spyAPI, () => true),
      ).toBeResolvedTo('SUCCESS');
      expect(spyAPI).toHaveBeenCalledTimes(1);
    });

    it('should retry with refreshed token when request fails', async () => {
      const currentTokenInfo = { currentAccessToken: 'TESTTOKEN' };
      const oauthAPI = createOAuthAPI(currentTokenInfo);
      let callCount = 0;
      const mockAPI = async (token: string) => {
        callCount++;
        if (token !== 'TESTTOKEN2') throw new Error('Invalid token');
        return 'SUCCESS';
      };
      await oauthAPI.requestAccessToken(false);
      currentTokenInfo.currentAccessToken = 'TESTTOKEN2';
      await expectAsync(
        oauthAPI.requestAPIWithOAuth(mockAPI, () => true),
      ).toBeResolvedTo('SUCCESS');
      expect(callCount).toBe(2);
    });

    it("should fail if refreshed token won't work", async () => {
      const currentTokenInfo = { currentAccessToken: 'TESTTOKEN' };
      const oauthAPI = createOAuthAPI(currentTokenInfo);
      let callCount = 0;
      const mockAPI = async () => {
        callCount++;
        throw new Error('Invalid token');
      };
      await oauthAPI.requestAccessToken(false);
      currentTokenInfo.currentAccessToken = 'TESTTOKEN2';
      await expectAsync(
        oauthAPI.requestAPIWithOAuth(mockAPI, () => true),
      ).toBeRejected();
      expect(callCount).toBe(2);
    });
  });
});
