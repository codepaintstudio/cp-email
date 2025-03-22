export const extractVariables = (content: string): string[] => {
  const regex = /\{(\w+)\}/g
  const matches = new Set<string>()
  let match: RegExpExecArray | null
  while ((match = regex.exec(content)) !== null) {
    matches.add(match[1])
  }
  return Array.from(matches)
}