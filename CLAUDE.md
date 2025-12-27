# DevNorth Backend - AI Assistant Guide

## Project Overview

This repository serves as the backend for the **DevNorth** project.

### Architecture & Technology Stack

- **Architecture Pattern**: Clean Architecture
- **Database**: PostgreSQL
- **Schema Management**: SQLC (SQL Compiler)
- **Database Migrations**: go-migrate
- **Language**: Go 1.25.2

## Working Philosophy

### 1. Ask Questions First

If you need more information to fulfill a task effectively, **ask clarifying questions before proceeding**. It's better to gather context upfront than to make incorrect assumptions.

### 2. Provide Honest Feedback

Critically evaluate user requests and ideas. Your value comes from thoughtful analysis, not blind agreement. Speak up when you identify:

- Flaws or potential issues in the proposed approach
- Opportunities for improvement
- Better alternatives or more efficient solutions
- Trade-offs the user should be aware of
- Technical debt or maintenance concerns

**Be direct and objective** - honest feedback saves time and improves outcomes.

### 3. POC Mindset

This is a **proof of concept**. Avoid over-engineering:

- Prioritize speed and simplicity over comprehensive solutions
- Not all features need production-grade implementation
- Focus on demonstrating core functionality
- Perfect is the enemy of done
- Ship working code over perfect code

### 4. Maintain Code Quality

Despite being a POC, **always adhere to software engineering principles**:

- Follow **SOLID principles**
- Write clean, readable, and maintainable code
- Use meaningful variable and function names
- Add comments only where logic isn't self-evident
- Ensure the codebase remains understandable and extensible
- Avoid creating technical debt that will be painful to fix later

**Quality ≠ Over-engineering** - Simple, clean code is quality code.

### 5. Navigate Trade-offs

When facing decisions between best practices and POC speed, **ask the user**:

- "Should I implement the full best-practice approach or take a faster, simplified route for this POC?"
- Present the trade-offs clearly:
  - **Option A**: Time/complexity vs. maintainability/scalability
  - **Option B**: Quick implementation vs. future technical debt
- Let the user make informed decisions based on project priorities
- Document the chosen approach and reasoning

### 6. Plan Before Acting

**Always follow this workflow:**

1. **Plan**: Design your approach and outline implementation steps
2. **Seek Approval**: Share your plan and ask for feedback or approval
3. **Implement Incrementally**: Execute the plan step-by-step
4. **Check-in Between Steps**: Request approval or feedback before proceeding to the next step
5. **Iterate**: Adjust based on feedback and learnings

**Never make large changes without user visibility and approval.**

## Clean Architecture Layers

When working with this codebase, respect the clean architecture boundaries:

```
┌─────────────────────────────────────┐
│         Delivery Layer              │  (HTTP handlers, gRPC, CLI)
│  (Frameworks & Drivers)             │
├─────────────────────────────────────┤
│         Use Cases Layer             │  (Business logic)
│  (Application Business Rules)       │
├─────────────────────────────────────┤
│         Domain Layer                │  (Entities, domain logic)
│  (Enterprise Business Rules)        │
├─────────────────────────────────────┤
│         Repository Layer            │  (Data access, SQLC queries)
│  (Interface Adapters)               │
└─────────────────────────────────────┘
```

**Dependency Rule**: Inner layers should never depend on outer layers. Dependencies point inward.

## Key Reminders

- **SQLC**: All database queries should be written in SQL and compiled using SQLC
- **go-migrate**: All schema changes must go through migrations
- **Testing**: When in doubt about test coverage, ask the user
- **Documentation**: Document decisions, especially when taking shortcuts for POC purposes
