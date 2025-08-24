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

```typescript
// * import structure example
import { onMount } from 'svelte'
import { useParams } from 'svelte-navigator'
import { UserIcon, Loader2Icon } from 'lucide-svelte'
import { Card, CardContent } from '$/lib/shadcn/components/ui/card'
import { backend, catcher } from '$/util/backend'
import type { PayloadProject } from '$/util/backend/backend.ts'
```

## Backend Integration Patterns

### API Client Usage

- **Backend Import**: Always use `backend` and `catcher` from `$/util/backend`
- **Promise Chain**: Use `.then()`, `.catch()`, `.finally()` pattern
- **Error Handling**: Use `catcher` utility for consistent error handling
- **Loading States**: Always handle loading states properly

```typescript
// * backend api call pattern
const loadProjectDetail = () => {
    loading = true
    backend.project
        .projectDetail({ projectId })
        .then((response) => {
            projectData = response.data.project
        })
        .catch((err) => {
            catcher(err)
        })
        .finally(() => {
            loading = false
        })
}
```

## UI & Styling Patterns

### Design System Integration

- **shadcn/ui Components**: Primary UI component library
    - Cards: `Card`, `CardContent`, `CardHeader`, `CardTitle`, `CardDescription`
    - Forms: `Input`, `Label`, `Button`, `Dialog` components
    - Tables: `Table`, `TableBody`, `TableCell`, `TableHead`, `TableHeader`, `TableRow`
    - Layout: Consistent spacing and typography
- **Tailwind CSS**: Utility-first styling with consistent design tokens
- **Lucide Icons**: Consistent iconography throughout the application
    - Use `size` prop for consistent icon sizing instead of `h-4 w-4`
    - Import specific icons to reduce bundle size

### Loading States Pattern

```typescript
// * loading indicator component
{#if loading}
    <div class="flex min-h-[400px] items-center justify-center">
        <Loader2Icon class="text-primary h-8 w-8 animate-spin" />
    </div>
{:else if !projectData}
    <div class="flex min-h-[400px] flex-col items-center justify-center">
        <InfoIcon class="mb-4 h-16 w-16 text-gray-400" />
        <h3 class="mb-2 text-lg font-semibold">Project not found</h3>
        <Button onclick={() => navigate('/project')}>Back to Projects</Button>
    </div>
{:else}
    <!-- * main content -->
{/if}
```

### Layout Patterns

- **Container Pattern**: Wrap page content in `<Container>` component
- **Grid Layouts**: Responsive grids (`grid-cols-1 md:grid-cols-2 lg:grid-cols-3`)
- **Flex Layouts**: Use flexbox for navigation and button groups
- **Fixed Positioning**: Floating components use `fixed bottom-center` positioning

### Card Design Pattern

```typescript
// * standard card layout with header and actions
<div class="flex items-center justify-between">
    <div class="flex items-center gap-3">
        <UserIcon class="h-6 w-6 text-purple-600" />
        <div>
            <h1 class="text-2xl font-bold">Users</h1>
            <p class="text-muted-foreground">Manage project members</p>
        </div>
    </div>
    <Button on:click={() => addUserDialogOpen = true}>
        <PlusIcon class="h-4 w-4" />
        Add User
    </Button>
</div>
```

## Component Development Patterns

### Dialog Management Pattern

```typescript
// * dialog state management
let addUserDialogOpen = false
let editUserDialogOpen = false
let editingUser: PayloadProjectUser | null = null

// * dialog event handlers
const handleUserAdded = () => {
    loadUsers()
    addUserDialogOpen = false
}

const handleEditUser = (user: PayloadProjectUser) => {
    editingUser = user
    editUserDialogOpen = true
}

// * dialog components at bottom of template
<AddUserDialog
    bind:open={addUserDialogOpen}
    {projectId}
    on:created={handleUserAdded}
/>

<EditUserDialog
    bind:open={editUserDialogOpen}
    {projectId}
    user={editingUser}
    on:updated={handleUserUpdated}
/>
```

### Table with Actions Pattern

```typescript
// * data table with action buttons
<Table>
    <TableHeader>
        <TableRow>
            <TableHead>User</TableHead>
            <TableHead>Email</TableHead>
            <TableHead>Role</TableHead>
            <TableHead class="text-right">Actions</TableHead>
        </TableRow>
    </TableHeader>
    <TableBody>
        {#each users as user (user.id)}
            <TableRow>
                <TableCell class="font-medium">{user.name}</TableCell>
                <TableCell class="text-muted-foreground">{user.email}</TableCell>
                <TableCell>
                    <Badge variant={getRoleBadgeVariant(user.role)}>
                        {user.role}
                    </Badge>
                </TableCell>
                <TableCell class="text-right">
                    <div class="flex justify-end gap-2">
                        <Button variant="ghost" size="sm" on:click={() => handleEditUser(user)}>
                            <EditIcon class="h-4 w-4" />
                        </Button>
                        <Button variant="ghost" size="sm" on:click={() => handleRemoveUser(user.id)}>
                            <TrashIcon class="h-4 w-4" />
                        </Button>
                    </div>
                </TableCell>
            </TableRow>
        {/each}
    </TableBody>
</Table>
```

### State Management

- **Local State**: Use `let` for component-specific state
- **Reactive Statements**: Use `$:` for computed values and side effects
- **Binding**: Two-way binding with `bind:` for form inputs

### Event Handling Patterns

- **Custom Events**: Dispatch typed events with `createEventDispatcher<T>`
- **Event Naming**: Use descriptive names (`created`, `updated`, `removed`)
- **Event Data**: Pass relevant data in event detail object
- **Handler Functions**: Separate functions for complex event logic

## Error Handling & User Experience

### Toast Notifications

```typescript
import { toast } from 'svelte-sonner'
toast.success('user added successfully')
toast.error(response.message)
```

## TypeScript Integration

### Type Safety

```typescript
// * import types explicitly
import type { PayloadProject, PayloadProjectUser } from '$/util/backend/backend.ts'

// * typed variables
let projectData: PayloadProject
let users: PayloadProjectUser[] = []
let editingUser: PayloadProjectUser | null = null
```

## Navigation & Routing

### Parameter Handling

```typescript
// path: src/page/project/[project]/user.svelte
export let project: number
```

## Key Points

- Always read `src/backend/backend.md` for backend types documentation
- Use `backend` and `catcher` from `$/util/backend` for all backend calls
- Implement proper loading states and user feedback
- Separate dialog components for better organization
- Use consistent naming patterns and file structure
- Handle edge cases like empty states and errors gracefully