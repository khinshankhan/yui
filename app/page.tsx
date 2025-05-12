import React from "react"

export default function Page() {
  return (
    <main className="flex min-h-screen items-center justify-center px-4">
      <article className="rounded-lg border-2 bg-slate-50 p-4 text-center">
        <h1 className="font-heading text-4xl font-semibold tracking-tight md:text-5xl lg:text-6xl">
          Yui
        </h1>
        <h2 className="font-body mt-6 text-lg font-medium tracking-normal md:text-xl lg:text-2xl">
          {">"} A bespoke kit of tools &mdash; each built to solve one small problem{" "}
          <em>very well</em>.
        </h2>
      </article>
    </main>
  )
}
