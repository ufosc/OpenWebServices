export const ValidateEmail = (email : string) => {
  if (email.match(/^[\w-\.]+@ufl.edu$/)) {
    return true
  }
  return false
}

export const ValidatePassword = (pwd : string) => {
  if (pwd.match(/^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{12,}$/)) {
    return true
  }
  return false
}

export const ValidateClientURLParams = (client : any) => {
  let hasDefined = false
  let hasUndefined = false
  const keys = ["response_type", "client_id", "redirect_uri", "state"]
  for (let i = 0; i < keys.length; i++) {
    if (client[keys[i]] === null) {
      hasUndefined = true
      continue
    }
    hasDefined = true
  }

  // User is either here on their own or they've been redirected by a client.
  // In the latter case, if any client parameter is defined, then all parameters
  // must be defined.
  if (hasDefined && hasUndefined) {
    return false
  }

  return true
}
