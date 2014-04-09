package invertedindex

import "path"

func docLocation(location string) string {
  return path.Join(location, "doc")
}

func docIndexLocation(location string) string {
  return path.Join(location, "docindex")
}

func indexLocation(location string) string {
  return path.Join(location, "index")
}

func termsLocation(location string) string {
  return path.Join(location, "terms")
}
