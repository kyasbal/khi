import { Pipe, PipeTransform } from '@angular/core';

export enum PrincipalType {
  System = 'System',
  Node = 'Node',
  ServiceAccount = 'SA',
  User = 'User',
  Invalid = 'Invalid',
  NotAvailable = 'N/A',
}

export interface ResourceOperatorPrincipal {
  type: PrincipalType;
  full: string;
  short: string;
}

/**
 * Parse the principal string modifying K8s resource into structured representation
 */
@Pipe({
  name: 'parsePrincipal',
})
export class ParsePrincipalPipe implements PipeTransform {
  transform(value: string): ResourceOperatorPrincipal {
    const result: ResourceOperatorPrincipal = {
      type: PrincipalType.User,
      full: value,
      short: value,
    };
    if (value === '') {
      result.type = PrincipalType.NotAvailable;
      result.full = '';
      result.short = '';
    }
    if (value.startsWith('system:serviceaccount:')) {
      result.type = PrincipalType.ServiceAccount;
      result.short = value.split('system:serviceaccount:')[1];
    } else if (value.startsWith('system:node:')) {
      result.type = PrincipalType.Node;
      result.short = value.split('system:node:')[1];
    } else if (value.startsWith('system:')) {
      result.type = PrincipalType.System;
      result.short = value.split('system:')[1];
    }
    return result;
  }
}
