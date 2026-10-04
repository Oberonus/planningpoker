// Restrict ticket links to web URLs because they become clickable in the game header.
export const ticketURLRules = [value => {
  if (!value || !value.trim()) {
    return true
  }

  try {
    const url = new URL(value.trim())
    return ['http:', 'https:'].includes(url.protocol) || 'Use an http:// or https:// link.'
  } catch (error) {
    return 'Enter a complete link, starting with https://.'
  }
}]
