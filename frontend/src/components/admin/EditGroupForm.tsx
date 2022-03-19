import { useEditGroup } from "../../hooks/admin/useEditGroup";
import { Group } from "../../model/Group";
import { GroupDetailsForm } from "./GroupDetailsForm";
interface GroupProps {
  group: Group;
}

export function EditGroupForm({ group }: GroupProps) {
  const mutation = useEditGroup(group.id);
  return (
    <GroupDetailsForm
      group={group}
      onSubmit={async (values) => {
        await mutation.mutateAsync(values);
      }}
    >
      Save Changes
    </GroupDetailsForm>
  );
}
