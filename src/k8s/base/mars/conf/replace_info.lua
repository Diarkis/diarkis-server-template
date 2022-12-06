function replace_sensitive_info(tag, timestamp, record)
  -- mask social security number
  record["log"] = string.gsub(record["log"], "%d%d%d%-+%d%d%-+%d%d%d%d", "xxx-xx-xxxx")
  -- mask credit card number
  record["log"] = string.gsub(record["log"], "%d%d%d%d *%d%d%d%d *%d%d%d%d *%d%d%d%d", "xxxx xxxx xxxx xxxx")
  -- mask email address
  record["log"] = string.gsub(record["log"], "[%w+%.%-_]+@[%w+%.%-_]+%.%a%a+", "user@email.tld")
  return 1, timestamp, record
end