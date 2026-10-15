import styled from "styled-components";

export const Container = styled.div`
    background-color: white;
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 2rem;
    gap: 1rem;
`
export const HeaderName = styled.div`
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
`
export const ContentBody = styled.div`
    flex: 1;
    display: flex;
    gap: 5rem;
`

export const ContentBorder = styled.div`
    display: flex;
    flex-direction: column;
    width: 100%;
`

export const ContentButtons = styled.div`
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    gap: 1rem;
`

export const MessageStatus = styled.div`
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background-color: #CECED2;
    border-radius: 1rem;

    > h1 {
        color: white;
    }
    > p {
        color: blue;
    }
`
