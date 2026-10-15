import { useState, useEffect } from "react";
import { Container, ContentBody, ContentBorder, ContentButtons, HeaderName, MessageStatus } from "./styles";
import { Button, Modal, Space, Table, message, Form, Input, Popconfirm } from "antd";
import SapatoService from "../../shared/service/Sapato/SapatoService";

export const Home = () => {
  const [modalOpen, setModalOpen] = useState(false);
  const [editingSapato, setEditingSapato] = useState(null);
  const [sapatos, setSapatos] = useState([]);
  const [statusMessage, setStatusMessage] = useState("");

  const [form] = Form.useForm();

  const handleGetSapatos = async () => {
    try {
      const response = await SapatoService.GetAllSapato();
      if (response.data && Array.isArray(response.data)) {
        setSapatos(response.data.map((sapato, index) => ({ ...sapato, key: index })));
        setStatusMessage("Cód: 200 - Success Loaded");
      } else {
        throw new Error("Dados inválidos recebidos");
      }
    } catch (error) {
      setStatusMessage(`Cód: ${error.response?.status || 500} - ${error.response?.data?.message || error.message}`);
      message.error("Erro ao carregar sapatos: " + error.message);
    }
  };

  const handleDeleteSapato = async (id) => {
    try {
      await SapatoService.DeleteSapato(id);
      setStatusMessage("Cód: 200 - Sapato Deletado");
      message.success("Sapato deletado com sucesso!");
      handleGetSapatos();
    } catch (error) {
      setStatusMessage(`Cód: ${error.response?.status || 500} - ${error.response?.data?.message || error.message}`);
      message.error("Erro ao deletar sapato: " + error.message);
    }
  };

  const handleGetHash =async () => {
    try {
      const response = await SapatoService.GetHash()
      if(response.data){
        setStatusMessage(response.data.hash)
      }
    } catch (error) {
      setStatusMessage(`Cód: ${error.response.status} - ${error.response.data?.message}`);
      message.error("Erro ao carregar o hash do arquivo: " + error.message);
    }
  }
  const handleDownload = async () => {
    try {
      const response = await SapatoService.Download();
      const blob = new Blob([response.data], { type: "application/zip" });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", "produtos.zip");
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      setStatusMessage("Cód: 200 - CSV Downloaded");
    } catch (error) {
      if(error.response){
        setStatusMessage(`Cód: ${error.response.status} - ${error.response.data?.message}`);
      }
      else{
        setStatusMessage('Erro no codigo');
      }
      message.error("Erro ao fazer download do CSV: " + error.message);
    }
  };

  const handleAddSapato = async (values) => {
    try {
      if (editingSapato) {
        await SapatoService.UpdateSapato(editingSapato.id, values);
        message.success("Sapato atualizado com sucesso!");
      } else {
        await SapatoService.CreateSapato(values);
        message.success("Sapato adicionado com sucesso!");
      }
      setModalOpen(false);
      form.resetFields();
      setEditingSapato(null);
      handleGetSapatos();
    } catch (error) {
      setStatusMessage(`Cód: ${error.response?.status || 500} - ${error.response?.data?.message || error.message}`);
      message.error("Erro ao salvar sapato: " + error.message);
    }
  };

  const handleEditSapato = (sapato) => {
    setEditingSapato(sapato);
    form.setFieldsValue(sapato);
    setModalOpen(true);
  };

  useEffect(() => {
    handleGetSapatos();
  }, []);

  const columns = [
    {
      title: "Marca",
      dataIndex: "marca",
      key: "marca",
    },
    {
      title: "Modelo",
      dataIndex: "modelo",
      key: "modelo",
    },
    {
      title: "Tamanho",
      dataIndex: "tamanho",
      key: "tamanho",
    },
    {
      title: "Cor",
      dataIndex: "cor",
      key: "cor",
    },
    {
      title: "Ações",
      key: "acoes",
      render: (_, record) => (
        <Space size="middle">
          <a onClick={() => handleEditSapato(record)}>Editar</a>
          <Popconfirm
            title="Tem certeza que deseja deletar este sapato?"
            onConfirm={() => handleDeleteSapato(record.id)}
            okText="Sim"
            cancelText="Não"
          >
            <a>Deletar</a>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <>
      <Container>
        <HeaderName>
          <h1>Sapataria Dois Irmões</h1>
        </HeaderName>
        <ContentBody>
          <ContentBorder>
            <Table dataSource={sapatos} columns={columns} pagination={{ pageSize: 6 }} />
          </ContentBorder>
          <ContentButtons>
            <Button type="primary" onClick={() => setModalOpen(true)}>
              Adicionar Sapato
            </Button>
            <Button type="primary" onClick={handleGetHash}>
              Gerar Hash
            </Button>
            <Button type="primary" onClick={handleDownload}>
              Download Csv
            </Button>
          </ContentButtons>
        </ContentBody>
        <MessageStatus>
          <h1>Msg Output:</h1>
          <p>{statusMessage}</p>
        </MessageStatus>
      </Container>

      <Modal
        title={editingSapato ? "Editar Sapato" : "Adicionar Sapato"}
        centered
        open={modalOpen}
        onOk={() => form.submit()}
        onCancel={() => {
          setModalOpen(false);
          form.resetFields();
          setEditingSapato(null);
        }}
      >
        <Form form={form} layout="vertical" onFinish={handleAddSapato}>
          <Form.Item
            label="Marca"
            name="marca"
            rules={[{ required: true, message: "Por favor, insira a marca do sapato!" }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            label="Modelo"
            name="modelo"
            rules={[{ required: true, message: "Por favor, insira o modelo do sapato!" }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            label="Tamanho"
            name="tamanho"
            rules={[{ required: true, message: "Por favor, insira o tamanho do sapato!" }]}
          >
            <Input type="number" />
          </Form.Item>
          <Form.Item
            label="Cor"
            name="cor"
            rules={[{ required: true, message: "Por favor, insira a cor do sapato!" }]}
          >
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};
