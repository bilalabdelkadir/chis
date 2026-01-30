interface PageHeaderProps {
  title: string;
  description?: string;
}

export function PageHeader({ title, description }: PageHeaderProps) {
  return (
    <header className="mb-8">
      <h1 className="text-xl font-semibold tracking-tight text-foreground sm:text-2xl">
        {title}
      </h1>
      {description && (
        <p className="mt-1.5 text-sm text-muted-foreground max-w-2xl">
          {description}
        </p>
      )}
    </header>
  );
}
