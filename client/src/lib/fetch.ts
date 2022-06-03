export function getCaller(
  baseUrl: string,
  basePath: string,
  service: string,
  token: string
) {
  return async function <Input, Output>(
    method: string,
    i: Input
  ): Promise<Output> {
    const headers: Record<string, string> = {
      Accept: "application/json",
      "Accept-Encoding": "gzip",
      "Content-Type": "application/json",
    };
    if (token != "") {
      headers["Authorization"] = `Bearer ${token}`;
    }
    const response = await fetch(
      `${baseUrl}/${basePath}/${service}.${method}`,
      {
        method: "POST",
        headers: headers,
        body: JSON.stringify(i),
      }
    );
    if (response.status != 200) {
      throw new Error(`bad status ${response.status} ${await response.text()}`);
    }
    return response.json();
  };
}

export function getToken(): string {
  return localStorage.getItem("access_token") || "";
}

export function setToken(token: string) {
  return localStorage.setItem("access_token", token);
}
