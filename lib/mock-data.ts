import { Note, Folder } from './types'

export const folders: Folder[] = [
  {
    id: '1',
    name: 'Học Go',
    icon: '📚',
    level: 0,
    children: [
      { id: '1-1', name: 'Cơ bản', icon: '📖', level: 1 },
      { id: '1-2', name: 'Goroutine', icon: '🔄', level: 1 },
    ],
  },
  {
    id: '2',
    name: 'Design Patterns',
    icon: '🏗️',
    level: 0,
    children: [
      { id: '2-1', name: 'Creational', icon: '⚙️', level: 1 },
      { id: '2-2', name: 'Structural', icon: '🔧', level: 1 },
    ],
  },
  {
    id: '3',
    name: 'Dự án cá nhân',
    icon: '💻',
    level: 0,
    children: [
      { id: '3-1', name: 'Note App', icon: '📝', level: 1 },
    ],
  },
]

export const notes: Note[] = [
  {
    id: '1',
    title: 'Introduction to Go',
    folder: '1-1',
    tags: ['golang', 'basics'],
    content: `# Introduction to Go

Go is a statically typed, compiled programming language designed at Google.

## Key Features
- Fast compilation
- Built-in concurrency
- Simple syntax

## Hello World
\`\`\`go
package main

import "fmt"

func main() {
  fmt.Println("Hello, World!")
}
\`\`\``,
    excerpt: 'Go is a statically typed, compiled programming language...',
    lastModified: new Date('2024-01-15'),
  },
  {
    id: '2',
    title: 'Understanding Goroutines',
    folder: '1-2',
    tags: ['golang', 'concurrency'],
    content: `# Understanding Goroutines

Goroutines are lightweight threads managed by the Go runtime.

## Creating Goroutines
\`\`\`go
go myFunction()
\`\`\`

## Key Concepts
- Lightweight compared to OS threads
- Managed by Go scheduler
- Communicate via channels`,
    excerpt: 'Goroutines are lightweight threads managed by the Go runtime...',
    lastModified: new Date('2024-01-14'),
  },
  {
    id: '3',
    title: 'Singleton Pattern',
    folder: '2-1',
    tags: ['design-pattern', 'creational'],
    content: `# Singleton Pattern

Ensures a class has only one instance and provides a global point of access.

## When to use
- Database connections
- Logger
- Configuration manager

## Implementation
Define a private constructor and provide a static method to get the instance.`,
    excerpt: 'Ensures a class has only one instance and provides a global point of access...',
    lastModified: new Date('2024-01-13'),
  },
  {
    id: '4',
    title: 'Adapter Pattern',
    folder: '2-2',
    tags: ['design-pattern', 'structural'],
    content: `# Adapter Pattern

Converts the interface of a class into another interface clients expect.

## Example
Making incompatible interfaces work together.

## Benefits
- Improves compatibility
- Promotes reusability
- Follows Open/Closed principle`,
    excerpt: 'Converts the interface of a class into another interface...',
    lastModified: new Date('2024-01-12'),
  },
  {
    id: '5',
    title: 'Note App Features',
    folder: '3-1',
    tags: ['note-app', 'project'],
    content: `# Note App Features

## Planned Features
1. Create and edit notes
2. Organize with folders
3. Tag system
4. Search functionality
5. Dark mode support

## Tech Stack
- React
- TypeScript
- Tailwind CSS
- Next.js`,
    excerpt: 'Planned features for the note application...',
    lastModified: new Date('2024-01-11'),
  },
  {
    id: '6',
    title: 'Go Concurrency Patterns',
    folder: '1-2',
    tags: ['golang', 'concurrency', 'patterns'],
    content: `# Go Concurrency Patterns

Common patterns for writing concurrent Go programs.

## Worker Pool Pattern
\`\`\`go
for i := 0; i < numWorkers; i++ {
  go worker(jobs, results)
}
\`\`\`

## Rate Limiting
Use time.Ticker for rate limiting goroutines.`,
    excerpt: 'Common patterns for writing concurrent Go programs...',
    lastModified: new Date('2024-01-10'),
  },
  {
    id: '7',
    title: 'Factory Pattern',
    folder: '2-1',
    tags: ['design-pattern', 'creational'],
    content: `# Factory Pattern

Provides an interface for creating families of related objects without specifying their concrete classes.

## Use Cases
- Object creation is complex
- Multiple implementations exist
- Construction logic should be centralized`,
    excerpt: 'Provides an interface for creating families of related objects...',
    lastModified: new Date('2024-01-09'),
  },
  {
    id: '8',
    title: 'API Design Best Practices',
    folder: '3-1',
    tags: ['project', 'api', 'best-practices'],
    content: `# API Design Best Practices

## RESTful Principles
- Use appropriate HTTP methods
- Meaningful URLs
- Proper status codes
- Versioning

## Documentation
Always provide clear API documentation for developers.`,
    excerpt: 'Best practices for designing scalable and maintainable APIs...',
    lastModified: new Date('2024-01-08'),
  },
]

export const allTags = Array.from(new Set(notes.flatMap(note => note.tags)))
