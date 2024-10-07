/**
 * metadata-types.ts
 * Defines the schema of metadata generated from KHI inspection task.
 * `Metadata` in kHI inspection task is a non main data(KHIFile) generated from a task run/dryrun.
 * Each inspection task generates the set of metadata and it is used in frontend.
 * Metadata includes the form fields needed to be filled on new inspection dialogs, task progress,..etc
 */

export type InspectionMetadataPlan = {
  taskGraph: string;
};

export type InspectionMetadataQuery = {
  id: string;
  name: string;
  query: string;
};

export type InspectionMetadataErrorSet = {
  errorMessages: InspectionMetadataError[];
};

export type InspectionMetadataError = {
  errorId: string;
  message: string;
  link: string;
};

export type InspectionMetadataFormFieldType = 'Text';

export type InspectionMetadataFormField = {
  id: string;
  type: InspectionMetadataFormFieldType;
  label: string;
  description: string;
  default: string;
  allowEdit: string;
  suggestions: string[];
  validationError: string;
  hint: string;
  hintType: 'warning' | 'info';
};

export type InspectionMetadataHeader = {
  inspectionType: string;
  inspectionTypeIconPath: string;
  inspectTimeUnixSeconds: number;
  startTimeUnixSeconds: number;
  endTimeUnixSeconds: number;
  suggestedFilename: string;
};

export type InspectionMetadataProgressPhase =
  | 'RUNNING'
  | 'ERROR'
  | 'CANCELLED'
  | 'DONE';

export type InspectionMetadataProgress = {
  phase: InspectionMetadataProgressPhase;
  progresses: InspectionMetadataProgressElement[];
  totalProgress: InspectionMetadataProgressElement;
};

export type InspectionMetadataProgressElement = {
  id: string;
  label: string;
  message: string;
  percentage: number;
  indeterminate: boolean;
};
