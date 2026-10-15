import { BrowserRouter } from 'react-router-dom'
import { Rotas } from './routes'
import { GlobalStyle } from './styles/globalStyles'

function App() {

  return (
    <>
    <GlobalStyle/>
      <BrowserRouter>
        <Rotas />
      </BrowserRouter>
    </>
  )
}

export default App
