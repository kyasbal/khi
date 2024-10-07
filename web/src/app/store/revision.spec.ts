import { RevisionState, RevisionVerb } from '../generated';
import { ResourceRevision } from './revision';

describe('ResourceRevision', () => {
  describe('get duration', () => {
    it('returns the calculated duration from startTime and endTime', () => {
      const revision = new ResourceRevision(
        1,
        10,
        RevisionState.RevisionStateExisting,
        RevisionVerb.RevisionVerbUpdate,
        '',
        '',
        false,
        false,
        0,
      );

      const duration = revision.duration;

      expect(duration).toBe(9);
    });
  });

  describe('get parsedManifest', () => {
    it('returns parsed YAML object from the resource content', () => {
      const revision = new ResourceRevision(
        0,
        1,
        RevisionState.RevisionStateExisting,
        RevisionVerb.RevisionVerbUpdate,
        `kind: foo`,
        '',
        false,
        false,
        0,
      );

      const manifest = revision.parsedManifest;

      expect(manifest?.kind).toBe('foo');
    });

    it('returns the cached YAML object', () => {
      const revision = new ResourceRevision(
        0,
        1,
        RevisionState.RevisionStateExisting,
        RevisionVerb.RevisionVerbUpdate,
        `kind: foo`,
        '',
        false,
        false,
        0,
      );

      const firstCall = revision.parsedManifest;
      const secoundCall = revision.parsedManifest;

      // The result must be an object with the same reference.
      expect(firstCall).toBe(secoundCall);
    });
  });
});
