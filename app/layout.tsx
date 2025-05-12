import React from "react"
import { cn } from "@/lib/utils"
import { GeistMono } from "geist/font/mono"
import { GeistSans } from "geist/font/sans"

import "./globals.css"

const fontAliases = {
  ["--font-heading"]: "var(--font-geist-sans)",
  ["--font-body"]: "var(--font-geist-sans)",
  ["--font-mono"]: "var(--font-geist-mono)",
} as React.CSSProperties

export default async function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={cn(GeistSans.variable, GeistMono.variable, "bg-emerald-400")}
      style={fontAliases}
    >
      <head></head>
      <body>{children}</body>
    </html>
  )
}
