const TRUE_IS_POSITIVE = new Set<string>(['Ready']);

export function isConditionPositive(
  resourceType: 'pod' | 'node',
  type: string,
  status: string,
): boolean {
  if (resourceType == 'pod') {
    return status == 'True';
  }
  if (TRUE_IS_POSITIVE.has(type)) {
    return status == 'True';
  } else {
    return status == 'False';
  }
}
