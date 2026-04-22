/**
 * Copyright 2026 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { ParserBlueprint } from 'src/app/parser/core/interfaces';
import { V6_BLUEPRINT } from 'src/app/parser/blueprints/v6-blueprint';

/**
 * Registry mapping file format versions to their specific parsing blueprints.
 *
 * Note: The KHI file format was originally documented as v7 in the design doc,
 * but the actual implementation uses v6. This is expected.
 */
export const VERSION_REGISTRY: Record<number, ParserBlueprint> = {
  6: V6_BLUEPRINT,
};
