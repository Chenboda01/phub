# phub Design System

## 0. Research Log

- Embedded refs: shortlisted `notion`, `warp`, and `linear.app`; picked `taste-skill` plus `notion` for an approachable, reading-first developer-tool page with a terminal-native focal object.
- UI/UX database: searched `terminal project launcher developer productivity minimal editorial`; retained the minimal single-column CTA discipline, responsive constraints, and accessible interaction guidance.
- Lazyweb: attempted the anonymous research route, but the MCP search request failed locally with curl error 43; no third-party screenshots were used or retained.
- Imagen drafts: skipped because no image-generation tool is available in this environment. The hero uses a semantic rendering of phub's actual prototype output instead of synthetic imagery.

## 1. Atmosphere & Identity

phub should feel like a tidy workbench: calm enough to read, precise enough to trust, and unmistakably terminal-native once the eye reaches the product specimen. The signature is the contrast between a quiet, moss-tinted document canvas and a deep charcoal terminal panel containing the actual first-milestone interface.

## 2. Color

### Palette

| Role | Token | Light | Dark | Usage |
| --- | --- | --- | --- | --- |
| Canvas | `--canvas` | `#f3f5f0` | `#101713` | Page background |
| Surface | `--surface` | `#ffffff` | `#19231d` | Navigation and content surfaces |
| Text | `--ink` | `#18221c` | `#f3f6f2` | Headlines and body copy |
| Muted text | `--muted` | `#59665f` | `#b2c0b8` | Supporting copy |
| Divider | `--line` | `#d8ded8` | `#304139` | Borders and rules |
| Accent | `--accent` | `#c8552a` | `#ef7b50` | Primary actions, links, focus |
| Accent pressed | `--accent-strong` | `#9b3d1f` | `#ff986e` | Active controls |
| Terminal | `--terminal` | `#142019` | `#0c120f` | Product specimen |
| Terminal text | `--terminal-ink` | `#eef5ef` | `#eef5ef` | Product specimen text |

### Rules

- The accent is reserved for interaction and the selected terminal row.
- Surfaces use tonal changes and whisper borders; no decorative gradients or glows.
- Light mode is the default; dark mode follows `prefers-color-scheme` with equivalent hierarchy.

## 3. Typography

| Level | Token | Size | Weight | Usage |
| --- | --- | --- | --- | --- |
| Display | `--type-display` | `clamp(3rem, 7vw, 5.75rem)` | 700 | Hero headline |
| H1 | `--type-heading` | `clamp(2rem, 4vw, 3.25rem)` | 700 | Section titles |
| H2 | `--type-subheading` | `1.5rem` | 700 | Panel titles |
| Body | `--type-body` | `1rem` | 400 | Reading copy |
| Small | `--type-small` | `0.875rem` | 500 | Labels and metadata |
| Mono | `--type-mono` | `0.875rem` | 500 | Commands and terminal output |

- Sans: `ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`.
- Mono: `ui-monospace, "SFMono-Regular", Consolas, "Liberation Mono", monospace`.
- External web fonts are intentionally avoided to preserve a fast, offline-friendly static page.
- CSS values: `--type-display: clamp(3rem, 7vw, 5.75rem)`, `--type-heading: clamp(2rem, 4vw, 3.25rem)`, `--type-subheading: 1.5rem`, `--type-body: 1rem`, `--type-small: 0.875rem`, and `--type-mono: 0.875rem`.

## 4. Spacing & Layout

- Base unit: 4px.
- Tokens: `--space-1` 4px, `--space-2` 8px, `--space-3` 12px, `--space-4` 16px, `--space-6` 24px, `--space-8` 32px, `--space-12` 48px, `--space-16` 64px, `--space-24` 96px.
- Content width: 1120px, with `clamp(20px, 5vw, 64px)` gutters.
- The hero is an asymmetric two-column layout at 900px and above, collapsing to one reading-first column below that breakpoint.
- Documentation sections are intentionally narrower than the product specimen to maintain a readable line length.
- CSS values: `--content-width: 70rem`, `--reading-width: 44rem`, `--gutter: clamp(1.25rem, 5vw, 4rem)`, and `--radius: 0.5rem`.

## 5. Components

### Masthead
- **Structure:** `header > nav > brand + anchor links + source link`.
- **States:** links underline and shift to accent on hover; every link has a visible accent focus ring.
- **Accessibility:** a skip link precedes the navigation; navigation landmarks and descriptive link labels are required.

### Action Link
- **Structure:** semantic anchor styled as a compact rounded rectangle.
- **Variants:** primary accent and quiet surface.
- **States:** hover darkens or lightens by theme, active translates 1px, focus uses a 3px accent outline.
- **Motion:** 180ms transform and color transition; disabled is not used on this static page.

### Terminal Specimen
- **Structure:** `figure > figcaption + pre` containing the current prototype's textual output.
- **States:** static; selection is communicated by the `>` character and accent color, never color alone.
- **Accessibility:** `aria-label` explains that it is current prototype output; line wrapping is prevented within the specimen but the container scrolls horizontally on narrow screens.

### Install Step
- **Structure:** ordinal label, concise explanation, and selectable code block.
- **States:** static documentation; links retain hover and focus states.
- **Accessibility:** commands remain text, not images; code blocks have an accessible label.

### Status Marker
- **Structure:** short text label with a small rule, not an icon-only badge.
- **States:** implemented and planned are named in text.
- **Accessibility:** status is never indicated by color alone.

## 6. Motion & Interaction

- Motion intensity: 3. The page uses only 180ms hover, active, and focus feedback to reinforce interaction.
- Any transition is limited to opacity, color, and transform.
- `prefers-reduced-motion: reduce` disables transitions and smooth scrolling.
- There are no auto-playing, scroll-triggered, or decorative animations.

## 7. Depth & Surface

- Strategy: mixed, using whisper borders for document structure and a single layered shadow for the terminal specimen.
- Content surfaces use `1px solid var(--line)` and an 8px radius.
- The terminal specimen uses a restrained three-layer tinted shadow to read as a physical terminal window without floating above the page.
- CSS values: `--radius-terminal: 0.75rem` and `--shadow-terminal: 0 1rem 2.5rem rgba(20, 32, 25, 0.18), 0 0.25rem 0.75rem rgba(20, 32, 25, 0.14), inset 0 1px 0 rgba(255, 255, 255, 0.08)`.

## 8. Accessibility Constraints & Accepted Debt

### Constraints

- Target: WCAG 2.2 AA, with visible keyboard focus, semantic landmarks, sequential headings, and a skip link.
- Body copy remains at least 16px; actions meet a 44px minimum target height.
- The page supports system dark mode and reduced motion.
- All primary actions work without JavaScript.

### Accepted Debt

| Item | Location | Why accepted | Owner / Exit |
| --- | --- | --- | --- |
| Static prototype specimen | Hero | It uses the exact current text output rather than a live embedded TUI, preserving a dependency-free Pages site. | Replace with a recorded artifact once a stable screenshot pipeline exists. |
