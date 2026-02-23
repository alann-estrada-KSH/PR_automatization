Act as an expert in {{.ProjectType}} applying the CDE method (Document as you code).

CDE means that the commit message IS living technical documentation of the system.
It is not a change log — it is the TECHNICAL RECORD of why the system evolved this way.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
STAGED FILES:
{{.Stats}}

DIFF:
```diff
{{.Diff}}
```

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ABSOLUTE RULES:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Reply ONLY with the commit message. No explanations, no extra text.
2. Follow EXACTLY the Conventional Commits format:

   <type>(<scope>): <imperative description, max 72 chars>
   <blank line>
   <body: 2-4 lines explaining the WHY, not the what>
   <blank line — ONLY if there is a footer>
   <optional footer: refs, BREAKING CHANGE>

3. The BODY is required when the change is non-trivial. Explain:
   - Why this technical decision was made
   - What problem existed before
   - Which function/class/module was the key point of the change
   If the change is trivial (typo, whitespace, lock file), body is optional.

4. Use the correct type:
   feat     → new functionality visible to the user or system
   fix      → bug correction with incorrect behavior
   refactor → code change without new features or bug fix (cleanup, reorganization)
   chore    → maintenance: deps, config, scripts, generated files
   docs     → documentation only (README, comments, docstrings)
   style    → formatting, whitespace, indentation (no logic change)
   test     → adding or correcting tests
   perf     → measurable performance improvement
   ci       → CI/CD configuration (GitHub Actions, Jenkins, etc.)
   build    → build system, Dockerfile, Makefile

5. <scope> is the affected module, component, or layer.
   Examples: auth, api, db, ui, payments, middleware, config, queue

6. Description: imperative present tense, lowercase start, no trailing period.
   ✅ "add token validation on refresh"
   ❌ "Added token validation." / "Adding token validation"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CDE COMMIT EXAMPLES (use these as quality reference):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

— Case: new feature with decision context
feat(auth): implement refresh token rotation in TokenService

Static access tokens represented an unlimited exposure window
if compromised. With rotation, each use of a refresh token invalidates
the previous one, limiting damage to a single session.
TokenService.refreshToken() now returns a new {access, refresh} pair.

— Case: bug fix with documented root cause
fix(payments): correct VAT calculation on orders with discounts

OrderCalculator.computeTax() was applying VAT on the gross total before
deducting promotions, resulting in overcharges.
Fixed operation order: discount → taxable base → VAT.

— Case: refactor with architecture rationale
refactor(api): extract permission validation to PermissionMiddleware

Permission logic was duplicated across 7 different controllers.
Centralizing into a middleware eliminates inconsistency and makes it
easier to audit which routes are restricted without reading each controller.

— Case: simple chore (no body required)
chore(deps): update security dependencies (February 2026)

— Case: small feat (clear scope, no body)
feat(ui): add skeleton loader on project listing

— Case: breaking change with footer
feat(api): require Authorization header on all endpoints

This version removes anonymous access to public routes deprecated
since v1.2. Clients must include a Bearer token on every request.

BREAKING CHANGE: /api/public/* routes no longer accept unauthenticated requests.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Generate the message for the diff shown above.
If the change is complex, write the body. If it is trivial, the subject line is enough.
