# Frontend Development Best Practices & Patterns

## Project Structure & Organization

### Component Architecture

- **Page Components**: Located in `src/page/` with feature-based directories
- **Shared Components**: Located in `src/component/` with type-based subdirectories
    - `ui/` - Reusable UI components (FloatingAction, PageTitle, Toggle)
    - `layout/` - Layout components (Container, Wrapper, AppLayout)
    - `screen/` - Screen-specific components (NotFound, WebviewNotice)
- **Dialog Components**: Separate dialogs in `page/[feature]/dialog/` for better organization
- **Floating Components**: Specialized floating UI in `component/ui/floating/`

### File Naming & Imports

- Use PascalCase for Svelte components (`.svelte` files)
- Import paths use `$` alias for src directory (`$/component/ui/...`)
- Group imports logically (Svelte, UI components, utilities, types)

## UI & Styling Patterns

### Design System Integration

- **shadcn/ui Components**: Primary UI component library
    - Cards: `Card`, `CardContent`, `CardHeader`, `CardTitle`, `CardDescription`
    - Forms: `Input`, `Label`, `Button`, `Dialog` components
    - Layout: Consistent spacing and typography
- **Tailwind CSS**: Utility-first styling with consistent design tokens
- **Lucide Icons**: Consistent iconography throughout the application
    - Use `size` prop for consistent icon sizing
    - Import specific icons to reduce bundle size

### Layout Patterns

- **Container Pattern**: Wrap page content in `<Container>` component
- **Grid Layouts**: Responsive grids (`grid-cols-1 md:grid-cols-2 lg:grid-cols-3`)
- **Flex Layouts**: Use flexbox for navigation and button groups
- **Fixed Positioning**: Floating components use `fixed bottom-center` positioning

### Card Design Pattern

- **Header Structure**: Icon + label + title layout
- **Content Structure**: Consistent spacing and typography
- **Hover Effects**: `transition-shadow hover:shadow-lg` for interactivity
- **Icon Integration**: Mini icons with consistent sizing and color

## Component Development Patterns

### Reusable Components

- **Props Interface**: Clear, typed props with defaults
- **Event Dispatching**: Use `createEventDispatcher` for parent communication
- **Slot Support**: Use `<slot />` for flexible content injection
- **Visibility Control**: `visible` prop for conditional rendering

### State Management

- **Local State**: Use `let` for component-specific state
- **Reactive Statements**: Use `$:` for computed values and side effects
- **Binding**: Two-way binding with `bind:` for form inputs

### Event Handling Patterns

- **Custom Events**: Dispatch typed events with `createEventDispatcher<T>`
- **Event Naming**: Use descriptive names (`search`, `paginate`, `created`)
- **Event Data**: Pass relevant data in event detail object
- **Handler Functions**: Separate functions for complex event logic

## Data & API Patterns

### Backend Integration

- **API Client**: Use generated `backend` client from OpenAPI spec
- **Error Handling**: Consistent error handling with `catcher` utility
- **Loading States**: Always handle loading states with proper UI feedback
- **Toast Notifications**: Use `svelte-sonner` for user feedback

### Pagination Implementation

- **Items Per Page**: Standardize pagination (24 items per page)
- **API Parameters**: Use `limit`, `offset`, `name` for consistent API calls
- **State Synchronization**: Keep pagination state in sync with API responses
- **Validation**: Validate page numbers and reset invalid inputs

### Search Functionality

- **Search Query**: Bind search input with debouncing consideration
- **Search Events**: Dispatch search events with query and page reset
- **Enter Key Support**: Handle Enter key for immediate search
- **Placeholder Text**: Descriptive placeholder text for context

## Form & Dialog Patterns

### Dialog Management

- **Separation of Concerns**: Move dialogs to separate components
- **State Management**: Handle open/close state and form data separately
- **Event Communication**: Use events for dialog-to-parent communication
- **Form Validation**: Client-side validation with user feedback
- **Loading States**: Disable forms during API calls

### Form Patterns

- **Input Binding**: Two-way binding for form inputs
- **Validation**: Real-time validation with error display
- **Submit Handling**: Async form submission with proper error handling
- **Reset Functionality**: Clear forms on success or cancel

## Animation & Transitions

### Svelte Transitions

- **Fly Animations**: Use `fly` transition for floating elements
- **Duration**: Consistent animation timing (200ms)
- **Direction**: Appropriate direction for element entry/exit

## Error Handling & User Experience

### Loading States

- **Loading Indicators**: Use `Loader2Icon` with spin animation
- **Empty States**: Meaningful empty state messages with actions
- **Error States**: Graceful error handling with user-friendly messages

### User Feedback

- **Toast Messages**: Success/error toasts for user actions
- **Button States**: Disabled states during loading
- **Visual Feedback**: Hover states and transitions for interactivity

## TypeScript Integration

### Type Safety

- **Import Types**: Import types explicitly (`import type { ... }`)
- **Event Types**: Type event dispatchers and handlers
- **API Types**: Use generated types from backend API
- **Component Props**: Type component props interfaces

## Performance Considerations

### Bundle Optimization

- **Selective Imports**: Import only needed components and utilities
- **Icon Optimization**: Import specific icons rather than entire libraries
- **Code Splitting**: Separate dialogs and heavy components

## Navigation & Routing

### Tab Navigation

- **Active States**: Visual indication of active tabs
- **Consistent Styling**: Shared styling patterns for tab interfaces
- **Icon Integration**: Icons with text labels for better UX

## Best Practices Summary

### Code Organization

- Separate concerns into focused components
- Use feature-based directory structure
- Group related functionality together

### User Interface

- Consistent design system usage
- Responsive design patterns
- Accessible component implementations

### Data Management

- Consistent API integration patterns
- Proper error handling and user feedback
- Efficient state management

### Component Design

- Reusable, composable components
- Clear prop interfaces and event contracts
- Proper separation of concerns
