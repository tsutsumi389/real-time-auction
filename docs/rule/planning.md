# Implementation Planning Guidelines

When creating implementation plans for new features or screens, follow these guidelines:

## Plan Document Format

**Location**: Store all implementation plans in `docs/plan/` directory

**File Naming**: Use descriptive names with underscores (e.g., `admin_login_implementation_plan.md`, `auction_create_implementation_plan.md`)

**Target Audience**: Write plans for non-technical stakeholders (project managers, designers, business analysts)
- Use plain Japanese with technical terms explained
- Avoid including actual code examples
- Focus on "what" needs to be built, not "how" to code it
- Use diagrams, tables, and flowcharts instead of code

## Required Sections

1. **Overview** (概要)
   - Purpose and goals
   - Target users
   - Key features

2. **Screen Layout** (画面レイアウト)
   - ASCII diagrams showing screen composition
   - Design requirements (colors, fonts, spacing)
   - Responsive design breakpoints

3. **Input Fields and Validation** (入力項目とバリデーション)
   - Table format with: label, input type, placeholder, required/optional, constraints
   - Validation rules with user-friendly error messages
   - Validation timing (onBlur, onSubmit, realtime)

4. **Processing Flow** (処理フロー)
   - Step-by-step flow diagrams
   - Success/failure scenarios
   - Error handling with user-facing messages

5. **Security Requirements** (セキュリティ要件)
   - Data handling (encryption, hashing)
   - Token management (JWT, sessions)
   - Common vulnerabilities addressed (SQL injection, XSS, CSRF)

6. **Database Design** (データベース設計)
   - Table structures in table format
   - Constraints and relationships
   - Test data specifications

7. **API Specification** (API仕様)
   - Endpoint paths and HTTP methods
   - Request/response JSON examples
   - Error response catalog

8. **Backend Implementation Requirements** (バックエンド実装要件)
   - Layer responsibilities (Domain, Repository, Service, Handler)
   - Describe "what each layer does" not "how to code it"
   - Middleware requirements
   - Technology stack

9. **Frontend Implementation Requirements** (フロントエンド実装要件)
   - Directory structure
   - Component descriptions (what they do)
   - State management requirements
   - Routing configuration

10. **Screen Transitions** (画面遷移)
    - Transition flow diagrams
    - Navigation guard behavior
    - Redirect rules

11. **UI Behavior Specification** (UIの動作仕様)
    - User interactions (click, keyboard, focus)
    - Loading states
    - Success/error feedback

12. **Responsive Design** (レスポンシブデザイン)
    - Mobile, tablet, desktop breakpoints
    - Layout adjustments per device

13. **Accessibility** (アクセシビリティ)
    - Keyboard navigation
    - Screen reader support
    - Focus management

14. **Test Requirements** (テスト要件)
    - Test scenarios (not test code)
    - Success criteria
    - Edge cases to cover

15. **Environment Variables** (環境変数)
    - Required configuration
    - Example values

16. **Implementation Steps** (実装手順)
    - Phase-by-phase breakdown with checkboxes (`- [ ]` format)
    - Time estimates per phase
    - Dependencies between phases
    - Total estimated time

17. **Success Criteria** (成功基準)
    - Checklist of completion requirements

18. **Next Steps** (次のステップ)
    - Related features to implement next

**Example**: See [docs/plan/admin_login_implementation_plan.md](../plan/admin_login_implementation_plan.md) for reference format

## What NOT to Include in Plans

- ❌ Actual Go/TypeScript/Vue code examples
- ❌ Implementation details (variable names, function signatures)
- ❌ Code snippets in any programming language (except JSON for API examples)
- ❌ Technical jargon without explanation

## What TO Include in Plans

- ✅ Plain language descriptions
- ✅ Diagrams and flowcharts
- ✅ Tables and structured data
- ✅ User-facing error messages
- ✅ Business logic flow
- ✅ JSON examples for API requests/responses
- ✅ SQL table definitions in table format
- ✅ Technology choices with rationale
