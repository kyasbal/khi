/**
 * api-types.ts
 * Defines the API schemas used between KHI backend and frontend
 */

import {
  InspectionMetadataErrorSet,
  InspectionMetadataFormField,
  InspectionMetadataHeader,
  InspectionMetadataPlan,
  InspectionMetadataProgress,
  InspectionMetadataQuery,
} from './metadata-types';

export interface InspectionType {
  id: string;
  name: string;
  description: string;
  icon: string;
}

export interface GetInspectionTypesResponse {
  types: InspectionType[];
}

export interface CreateInspectionTaskResponse {
  inspectionId: string;
}

export interface InspectionFeature {
  id: string;
  label: string;
  description: string;
  enabled: boolean;
}

export type InspectionDryRunResponse = {
  metadata: InspectionMetadataInDryrun;
};

export type InspectionDryRunRequest = InspectionRunRequest;
export type InspectionRunRequest = { [key: string]: unknown };

export type InspectionMetadataLog = {
  id: string;
  name: string;
  log: string;
};

export type InspectionMetadataResponse = InspectionMetadataOfRunResult;

export type InspectionMetadataInDryrun = {
  form: InspectionMetadataFormField[];
  query: InspectionMetadataQuery[];
  plan: InspectionMetadataPlan;
};

export type InspectionMetadataInTaskList = {
  progress: InspectionMetadataProgress;
  header: InspectionMetadataHeader;
  error: InspectionMetadataErrorSet;
};

export type InspectionMetadataOfRunResult = {
  header: InspectionMetadataHeader;
  query: InspectionMetadataQuery[];
  plan: InspectionMetadataPlan;
  log: InspectionMetadataLog[];
  error: InspectionMetadataErrorSet;
};

export type GetInspectionTasksResponse = {
  tasks: {
    [taskId: string]: InspectionMetadataInTaskList;
  };
  serverStat: {
    totalMemoryAvailable: number;
  };
};

export interface GetInspectionTaskFeatureResponse {
  features: InspectionFeature[];
}

export interface PutInspectionTaskFeatureRequest {
  features: string[];
}

export type PopupFormType = 'text' | 'popup_redirect';

/**
 * PopupFormRequest is a type returned on the endpoint GET /api/v2/popup.
 * Note this request is from backend with polling. Thus this is also a response in HTTP.
 */
export interface PopupFormRequest {
  id: string;
  title: string;
  type: PopupFormType;
  description: string;
  placeholder: string;
  options: {
    /**
     * The redirect target. This option is valid only when the type is `popup_redirect`.
     */
    redirectTo?: string;
    [key: string]: string | undefined;
  };
}

/**
 * PopupAnswerResponse is a type replied to the server with the endpoint POST /api/v2/popup/answer or POST /api/v2/popup/validate
 */
export interface PopupAnswerResponse {
  id: string;
  value: string;
}

/**
 * PopupValidationResult is a type returned from server on POST /api/v2/popup/validate
 */
export interface PopupAnswerValidationResult {
  id: string;
  validationError: string;
}
