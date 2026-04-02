---
trigger: glob
globs: **/*.ts, **/*.html, **/*.scss
---

# KHI Angular Standards

When developing or modifying Angular code in the KHI project, you **must** adhere to the following rules and best practices.

## 1. Verifications

1. **Build Verification**: Run `make build-web` to check compilation.
2. **Lint Verification**: Run `make lint-web` to check style.
3. **Test Verification**: Run `make test-web-headless` to check unit tests.
4. **Storybook Verification**: Run `make build-storybook` to check compilation for storybook.
5. **Review Verification**: Before asking the user to verify your changes, you MUST invoke a subagent for code review. See Section 4 for the procedure. **If you make any modifications based on the review, you MUST run the review again to verify the changes.**
6. **Restart on Correction**: If you make any corrections during a verification phase, you MUST restart the verification from the beginning for that phase.

## 2. General Coding Rules

1. **Comments**:
   - Use TSDoc-style comments for all public types, functions, and methods.
2. **Naming Conventions**:
   - Component selectors should have a `khi-` prefix.
3. **Test File Naming**: Test files must be named with `A.spec.ts` if `A.ts` exists. Do not define spec files by scenarios.
4. **Type Safety**: Do not use `any`.

## 3. Language-Specific Conventions (TypeScript & SCSS)

### 3.1 Modern Angular Conventions

These rules apply when creating new components or refactoring existing ones:

1. **Standalone Components**: Explicitly list dependencies in the `imports` array.
2. **Signals Paradigm**:
   - Use `input()` or `input.required()` for component inputs instead of `@Input`.
   - Use `output()` instead of `@Output`.
   - Use `model()` for two-way bindings.
   - Use `signal()` and `computed()` for component-level state management.
3. **Control Flow**: Use built-in control flow (`@if`, `@for`, `@switch`) in templates instead of structural directives (`*ngIf`, `*ngFor`).
4. **Access Modifiers**:
   - Members accessed only from the template must be `protected`.
   - Members used only within the TypeScript class must be `private`.
   - Use `readonly` for properties that are not reassigned.
5. **RxJS to Signal Conversion**: If a service returns an Observable, convert it to a Signal in the component using `toSignal`.
6. **Icons**: When importing `MatIconModule`, you must also import `KHIIconRegistrationModule`.
7. **File Separation**: Styles and templates must be defined in separate files. Do not directly supply them in `@Component`.

### 3.2 General Coding Rules for SCSS

1. **No Color Literals**: Do not use color literals in SCSS files, such as `background-color: #FF00FF;`. Define semantically meaningful color variables at the top of the SCSS file.
2. **Color Palette**: Use color palette from Material with `mat.m2-get-color-from-palette` rather than specifying color codes.
3. **Layout**: Prefer `display: grid` over `display: flex`. Use the `grid-template` property rather than specifying `grid-template-areas`, `grid-template-columns`, or `grid-template-rows` separately.
4. **Theme**: KHI's color scheme uses a light theme.

## 4. Subagent Review Guidelines

> [!IMPORTANT]
> **DO NOT FORGET** to invoke the subagent for code review after making changes. You must complete the subagent review before asking the user to verify your implementation.

Follow these rules to perform code reviews using a temporary subagent before asking the user to verify your changes.

- **Define**: Use `define_subagent` to create a temporary subagent specialized for Angular and TypeScript.
- **Invoke**: Pass the modified file paths (TS, HTML, SCSS) to the subagent during `invoke_subagent`.
- **Capabilities**: The subagent can read files and perform web searches.
- **Timeout and Retry**: You MUST set set 180sec as a deadline duration. If a subagent does not respond within the expected time or seems to be stuck, invoke a new subagent to retry the task.

### Review Checklist Focus

The subagent must verify:

- **Simplicity & Duplication**: Can any parts be written more simply and concisely? Are there any duplicated implementations?
- **Signals & Modern APIs**: Are Signals (`input`, `output`, `computed`, etc.) and built-in control flow (`@if`, `@for`) used correctly instead of legacy decorators and directives?
- **Access Modifiers**: Are members accessed only from templates marked as `protected`? Are others properly scoped (`private` or `public` readonly)?
- **Style & SCSS**: Are color literals avoided in SCSS? Are semantic variables or Material palettes used?
- **Standalone Components**: Are all component dependencies listed in the `imports` array? Is `KHIIconRegistrationModule` included if `MatIconModule` is present?
- **Testing**: Do test cases sufficiently cover realistic component interactions and state changes?

### Example Format

**`define_subagent`**

```json
{
  "name": "temp_angular_reviewer",
  "description": "Reviews Angular code against project standards.",
  "prompt_sections": [
    {
      "title": "Checklist",
      "content": "- Code simplicity and no duplication\n- Use of Signals and modern control flow\n- Proper access modifiers (protected for template accessed)\n- No color literals in SCSS\n- Independent style and template files\n- Proper icon module registration"
    }
  ],
  "tool_names": ["view_file", "search_web"]
}
```

**`invoke_subagent`**

```json
{
  "TypeName": "temp_angular_reviewer",
  "Role": "Angular Code Reviewer",
  "Prompt": "Review the changes in: [file paths]"
}
```
