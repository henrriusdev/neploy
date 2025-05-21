import { ThemeSwitcher } from "@/components/theme-switcher";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { useIsMobile, useTheme } from "@/hooks";
import { cn } from "@/lib/utils";
import { Menu } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import ReactMarkdown from "react-markdown";

interface Heading {
  id: string;
  text: string;
  level: number;
}

export default function Manual({ content }: { content: string }) {
  const [headings, setHeadings] = useState<Heading[]>([]);
  const [activeId, setActiveId] = useState<string>("");
  const contentRef = useRef<HTMLDivElement>(null);
  const headingRefs = useRef<Record<string, HTMLElement>>({});
  const isMobile = useIsMobile();
  const { theme, isDark, applyTheme } = useTheme();

  useEffect(() => {
    applyTheme(theme, isDark);
  }, [theme, isDark, applyTheme]);

  // Extract headings from markdown content
  useEffect(() => {
    const extractHeadings = (markdown: string) => {
      const headingRegex = /^(## |### )(.+)$/gm;
      const extractedHeadings: Heading[] = [];
      let match;

      while ((match = headingRegex.exec(markdown)) !== null) {
        const level = match[1].trim() === "##" ? 2 : 3;
        const text = match[2].trim();
        const id = text
          .toLowerCase()
          .replace(/\s+/g, "-")
          .replace(/[^\w-]/g, "");

        extractedHeadings.push({ id, text, level });
      }

      return extractedHeadings;
    };

    setHeadings(extractHeadings(content));
  }, [content]);

  // Set up intersection observer to track active heading
  useEffect(() => {
    if (!contentRef.current) return;

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setActiveId(entry.target.id);
          }
        });
      },
      {
        rootMargin: "0px 0px -80% 0px",
        threshold: 0.1,
      }
    );

    // Observe all heading elements
    const elements = contentRef.current.querySelectorAll("h2, h3");
    elements.forEach((element) => {
      headingRefs.current[element.id] = element as HTMLElement;
      observer.observe(element);
    });

    return () => {
      elements.forEach((element) => observer.unobserve(element));
    };
  }, [content, headings]);

  // Scroll to heading when clicked in TOC
  const scrollToHeading = (id: string) => {
    const element = headingRefs.current[id];
    if (element) {
      window.scrollTo({
        top: element.offsetTop - 100,
        behavior: "smooth",
      });
    }
  };

  // Custom components for ReactMarkdown
  const components = {
    h1: ({ node, ...props }: any) => {
      const id = props.children
        .toString()
        .toLowerCase()
        .replace(/\s+/g, "-")
        .replace(/[^\w-]/g, "");

      return (
        <h1
          id={id}
          className="scroll-mt-20 text-primary text-2xl lg:text-4xl font-bold"
          {...props}
        />
      );
    },
    h2: ({ node, ...props }: any) => {
      const id = props.children
        .toString()
        .toLowerCase()
        .replace(/\s+/g, "-")
        .replace(/[^\w-]/g, "");

      return <h2 id={id} className="scroll-mt-20 text-black font-semibold dark:text-accent text-xl lg:text-2xl" {...props} />;
    },
    h3: ({ node, ...props }: any) => {
      const id = props.children
        .toString()
        .toLowerCase()
        .replace(/\s+/g, "-")
        .replace(/[^\w-]/g, "");

      return <h3 id={id} className="scroll-mt-20 text-secondary text-lg lg:text-xl" {...props} />;
    },
  };

  // Table of Contents component
  const TableOfContents = () => (
    <div className="w-full">
      <ThemeSwitcher className="mb-4" />
      <h3 className="mb-4 text-lg font-semibold">Table of Contents</h3>
      <ul className="space-y-1">
        {headings.map((heading) => (
          <li key={heading.id}>
            <Button
              variant="ghost"
              className={cn(
                "w-full justify-start px-2 text-left",
                heading.level === 3 && "pl-6",
                activeId === heading.id
                  ? "bg-primary/10 text-primary font-medium"
                  : "text-muted-foreground hover:text-foreground"
              )}
              onClick={() => scrollToHeading(heading.id)}>
              {heading.text}
            </Button>
          </li>
        ))}
      </ul>
    </div>
  );

  return (
    <div className="flex flex-col md:flex-row w-full min-h-screen">
      {/* Mobile sidebar */}
      {isMobile && (
        <Sheet>
          <SheetTrigger asChild>
            <Button
              variant="outline"
              size="icon"
              className="fixed top-4 left-4 z-40">
              <Menu className="h-5 w-5" />
              <span className="sr-only">Toggle table of contents</span>
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-[280px] sm:w-[350px]">
            <ScrollArea className="h-[calc(100vh-4rem)] py-4">
              <TableOfContents />
            </ScrollArea>
          </SheetContent>
        </Sheet>
      )}

      {/* Desktop sidebar */}
      {!isMobile && (
        <div className="hidden md:block w-64 lg:w-72 shrink-0 h-screen sticky top-0 border-r border-border">
          <ScrollArea className="h-screen py-8 px-4">
            <TableOfContents />
          </ScrollArea>
        </div>
      )}

      {/* Main content */}
      <div
        ref={contentRef}
        className="flex-1 px-4 md:px-8 py-12 max-w-3xl mx-auto">
        <div className="prose prose-slate dark:prose-invert max-w-none">
          <ReactMarkdown components={components}>{content}</ReactMarkdown>
        </div>
      </div>
    </div>
  );
}
