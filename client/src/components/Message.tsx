import Column from "../dawn-ui/components/Column";
import Icon from "../dawn-ui/components/Icon";
import Row from "../dawn-ui/components/Row";
import Message from "../syrenity-client/structures/Message";

export default function MessageC({message}: {message: Message}) {
  return (
    <Row util={["no-shrink"]}>
      <Icon size="48px" src="" fallback="/images/logos/no_shape_logo.png" />
      <Column util={["flex-grow"]} style={{gap:"5px"}}>
        <Row util={["align-center"]} style={{gap: "10px"}}>
          <b>{message.authorId}</b>
          <small>{message.createdAt.toLocaleString()}</small>
        </Row>
        <label>{message.content}</label>
      </Column>
    </Row>
  );
}