export interface Theme {
  name: string;
  firstColor: string;
  secondColor: string;
  thirdColor: string;
  fourthColor: string;
  textColor: string;
  filter: string;
}

const darkTheme: Theme = {
  name: "Dark",
  firstColor: "#000000",
  secondColor: "rgba(7, 7, 7, 0.93)",
  thirdColor: "rgba(19, 19, 19, 0.9)",
  fourthColor: "rgba(105, 102, 102, 0.93)",
  textColor: "white",
  filter: "invert()",
};

const mintTheme: Theme = {
  name: "Mint",
  firstColor: "#8EEEB3",
  secondColor: "#FFFFFF",
  thirdColor: "#FFFFFF",
  fourthColor: "#b2a198",
  textColor: "black",
  filter: "invert(30%)",
};

const asiimov: Theme = {
  name: "Asiimov",
  firstColor: "#F67E23E5",
  secondColor: "white",
  thirdColor: "#ffaa69",
  fourthColor: "#fff8f4",
  textColor: "black",
  filter: "none",
};

export const AllThemes: Array<Theme> = [mintTheme, darkTheme, asiimov];
export const AllThemeNames: Array<string> = AllThemes.map((t) => t.name);

export function GetThemeInterface(th: string) {
  const i = AllThemeNames.indexOf(th);
  if (i < 0) {
    return { ...darkTheme };
  }
  return { ...AllThemes[i] };
}
