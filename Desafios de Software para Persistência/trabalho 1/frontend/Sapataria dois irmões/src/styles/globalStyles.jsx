import { createGlobalStyle } from "styled-components";

export const GlobalStyle = createGlobalStyle`   
    @font-face {
        font-family: 'Kantumruy Pro';
        src: url('https://fonts.googleapis.com/css2?family=Kantumruy+Pro:ital,wght@0,100..700;1,100..700&family=Poppins:wght@400;500;700&display=swap')
    }

    * {
        font-family: 'Kantumruy Pro', sans-serif;
        margin: 0;
        padding: 0;
        box-sizing: border-box;
        outline: 0;
        
  button {
    cursor: pointer;
    }

    }

    body {
        background-color: white;
    }

    .MuiOutlinedInput-root .MuiOutlinedInput-notchedOutline {
    border: none;
    }
`