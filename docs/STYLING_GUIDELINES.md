# P2K16 Design System & Styling Guidelines

## Overview

This document defines the design system for the P2K16 hackerspace management application, ensuring consistency and usability across all interfaces.

## Design Principles

### 1. **Accessibility First**
- WCAG 2.1 AA compliance
- Keyboard navigation support
- Screen reader friendly
- High contrast options

### 2. **Mobile-First Responsive Design**
- Progressive enhancement
- Touch-friendly interfaces
- Fluid layouts that work on all screen sizes

### 3. **Clear Information Hierarchy**
- Logical content organization
- Consistent spacing and typography
- Visual emphasis on important actions

### 4. **Feedback and State Communication**
- Clear loading states
- Informative error messages
- Success confirmations
- System status indicators

## Color Palette

### Primary Colors
```css
/* Primary Blue - Used for main actions and highlights */
--primary-blue: #2563eb;
--primary-blue-hover: #1d4ed8;
--primary-blue-light: #dbeafe;

/* Secondary Colors */
--secondary-gray: #6b7280;
--secondary-gray-light: #f3f4f6;
--secondary-gray-dark: #374151;
```

### Semantic Colors
```css
/* Success - For positive actions and confirmations */
--success-green: #10b981;
--success-green-light: #d1fae5;

/* Warning - For caution and pending states */
--warning-yellow: #f59e0b;
--warning-yellow-light: #fef3c7;

/* Danger - For errors and destructive actions */
--danger-red: #ef4444;
--danger-red-light: #fee2e2;

/* Info - For informational messages */
--info-blue: #3b82f6;
--info-blue-light: #dbeafe;
```

### Neutral Colors
```css
/* Text colors */
--text-primary: #111827;
--text-secondary: #6b7280;
--text-muted: #9ca3af;

/* Background colors */
--bg-primary: #ffffff;
--bg-secondary: #f9fafb;
--bg-tertiary: #f3f4f6;

/* Border colors */
--border-light: #e5e7eb;
--border-medium: #d1d5db;
--border-dark: #9ca3af;
```

## Typography

### Font Stack
```css
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
```

### Type Scale
```css
/* Headings */
--text-4xl: 2.25rem;  /* 36px - Page titles */
--text-3xl: 1.875rem; /* 30px - Section headers */
--text-2xl: 1.5rem;   /* 24px - Component titles */
--text-xl: 1.25rem;   /* 20px - Subheadings */
--text-lg: 1.125rem;  /* 18px - Large body text */

/* Body text */
--text-base: 1rem;    /* 16px - Default body text */
--text-sm: 0.875rem;  /* 14px - Small text */
--text-xs: 0.75rem;   /* 12px - Captions, labels */
```

### Font Weights
```css
--font-weight-light: 300;
--font-weight-normal: 400;
--font-weight-medium: 500;
--font-weight-semibold: 600;
--font-weight-bold: 700;
```

## Spacing Scale

Based on 8px baseline grid:

```css
--space-1: 0.25rem;  /* 4px */
--space-2: 0.5rem;   /* 8px */
--space-3: 0.75rem;  /* 12px */
--space-4: 1rem;     /* 16px */
--space-5: 1.25rem;  /* 20px */
--space-6: 1.5rem;   /* 24px */
--space-8: 2rem;     /* 32px */
--space-10: 2.5rem;  /* 40px */
--space-12: 3rem;    /* 48px */
--space-16: 4rem;    /* 64px */
--space-20: 5rem;    /* 80px */
```

## Component Guidelines

### Buttons

#### Primary Buttons
- Used for main actions (login, save, create)
- Blue background with white text
- Rounded corners (4px)
- Adequate padding for touch targets (min 44px height)

#### Secondary Buttons
- Used for secondary actions (cancel, back)
- Outlined style with primary color
- Same dimensions as primary buttons

#### Danger Buttons
- Used for destructive actions (delete, logout)
- Red color scheme
- Clear confirmation required

### Form Elements

#### Input Fields
- Consistent border radius (4px)
- Clear focus states with color and shadow
- Proper spacing between fields (16px)
- Labels positioned above inputs
- Placeholder text in muted color

#### Validation
- Inline validation for immediate feedback
- Error states with red color and clear messaging
- Success states with green indicators

### Cards and Containers

#### Cards
- Subtle shadow for depth
- Rounded corners (8px)
- Consistent padding (16px)
- White background on light themes

#### Spacing
- Consistent margins between sections (24px)
- Proper content hierarchy with spacing
- Adequate whitespace for readability

### Navigation

#### Header Navigation
- Fixed header with site branding
- User account dropdown in top right
- Clear active state indicators
- Mobile-responsive hamburger menu

#### Breadcrumbs
- For complex navigation hierarchies
- Clear path indication
- Clickable parent levels

### Status and Feedback

#### Alerts
- Distinct styling for different message types
- Dismissible when appropriate
- Clear icons for message type
- Proper contrast ratios

#### Loading States
- Skeleton screens for content loading
- Spinner for action feedback
- Progress indicators for long operations

#### Empty States
- Helpful messaging for empty data
- Clear calls to action
- Illustrative content when helpful

## Responsive Design

### Breakpoints
```css
/* Mobile first approach */
--breakpoint-sm: 640px;   /* Small tablets */
--breakpoint-md: 768px;   /* Tablets */
--breakpoint-lg: 1024px;  /* Laptops */
--breakpoint-xl: 1280px;  /* Desktops */
```

### Grid System
- 12-column grid system
- Flexible container widths
- Consistent gutter spacing (16px)

## Accessibility Guidelines

### Color Contrast
- Minimum 4.5:1 contrast ratio for normal text
- Minimum 3:1 contrast ratio for large text
- Color should not be the only way to convey information

### Focus Management
- Visible focus indicators
- Logical tab order
- Focus trapping in modals
- Skip links for main content

### Screen Reader Support
- Semantic HTML elements
- Proper heading hierarchy
- Alt text for images
- ARIA labels where needed

## Implementation Notes

### CSS Custom Properties
All design tokens should be implemented as CSS custom properties for easy theming and maintenance.

### Component Classes
Use BEM methodology for CSS class naming:
- Block: `.p2k16-button`
- Element: `.p2k16-button__icon`
- Modifier: `.p2k16-button--primary`

### HTMX Integration
- Use consistent classes for HTMX targets
- Implement smooth transitions for content updates
- Provide loading states during requests

## File Organization

```
/styles/
├── tokens/           # Design tokens (colors, spacing, etc.)
├── base/            # Reset, typography, base styles
├── components/      # Component-specific styles
├── utilities/       # Utility classes
└── main.css        # Main stylesheet combining all modules
```

## Usage Examples

### Button Implementation
```html
<button class="p2k16-button p2k16-button--primary">
  <span class="p2k16-button__text">Save Changes</span>
</button>
```

### Card Component
```html
<div class="p2k16-card">
  <div class="p2k16-card__header">
    <h3 class="p2k16-card__title">User Profile</h3>
  </div>
  <div class="p2k16-card__body">
    <!-- Card content -->
  </div>
</div>
```

### Form Field
```html
<div class="p2k16-field">
  <label class="p2k16-field__label" for="username">Username</label>
  <input class="p2k16-field__input" type="text" id="username" name="username">
  <div class="p2k16-field__error" id="username-error"></div>
</div>
```

## Testing Checklist

- [ ] Test with screen readers
- [ ] Verify keyboard navigation
- [ ] Check color contrast ratios
- [ ] Test on various screen sizes
- [ ] Validate with accessibility tools
- [ ] Cross-browser compatibility testing

---

*This design system should be reviewed and updated regularly to ensure it continues to meet user needs and accessibility standards.*