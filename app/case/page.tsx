"use client"

import React, { useState } from "react"

interface CaseBoxProps {
  input: string
  name: string
  transformer: (str: string) => string
}

function CaseBox({ input, name, transformer }: CaseBoxProps) {
  return (
    <li className="w-full rounded-md bg-slate-50 p-4 text-left shadow-lg">
      <label className="mb-1 block text-sm font-semibold text-gray-700">{name}</label>

      <span className="inline-block w-full text-center text-lg">{transformer(input)}</span>
    </li>
  )
}

const casings: {
  name: string
  transformer: (str: string) => string
}[] = [
  {
    name: "Upper Case",
    transformer: (str: string) => str.toUpperCase(),
  },
  {
    name: "Lower Case",
    transformer: (str) => str.toLowerCase(),
  },
  {
    name: "Kebab Case",
    transformer: (str) => str.split(" ").join("-"),
  },
  {
    name: "Snake Case",
    transformer: (str) => str.split(" ").join("_"),
  },
  {
    name: "Camel Case",
    transformer: (str) =>
      str
        .split(" ")
        .map((w) => w.charAt(0).toLowerCase() + w.slice(1))
        .join(""),
  },
  {
    name: "Pascal Case",
    transformer: (str) =>
      str
        .split(" ")
        .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
        .join(""),
  },
]

export default function Page() {
  const [input, setInput] = useState("")

  return (
    <main className="flex min-h-screen flex-col items-center px-4 py-24">
      <div className="w-full max-w-xl text-center">
        <h1 className="text-3xl font-bold">Case Converter</h1>
        <input
          type="text"
          placeholder="Enter a string..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          className="mt-8 w-4/5 rounded-sm border bg-slate-50 p-2.5 text-lg"
        />
        <ul className="mt-10 flex flex-col gap-6">
          {casings.map((casing) => {
            return (
              <CaseBox
                key={casing.name}
                input={input}
                name={casing.name}
                transformer={casing.transformer}
              />
            )
          })}
        </ul>
      </div>
    </main>
  )
}
