export function snake_case_keys(obj: any): string[] {
    return  Object.keys(obj).map((key) => key.replace(/\W+/g, " ")
    .split(/ |\B(?=[A-Z])/)
    .map(word => word.toLowerCase())
    .join('_'));
}

export function all_keys(obj: any): string[] {
    return  Object.keys(obj).map((key) => key);
}